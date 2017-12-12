package SourceProducer

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/YiftachE/DockerHPC/Core"
	"github.com/fsnotify/fsnotify"
)

type FileProducer struct {
	config     FileProducerConfig
	rabbitConn Core.RabbitConnection
}

func (producer FileProducer) NotifyOnClose(error chan error) {

}
func CreateProducer(path string) (FileProducer, error) {
	config, err := GetConfig(path)
	if err != nil {
		return FileProducer{}, err
	}
	conn, connectionError := Core.NewRabbitConnection(config.RabbitConnectionString)
	if connectionError != nil {
		return FileProducer{}, connectionError
	}
	return FileProducer{config, conn}, nil
}

func (producer FileProducer) Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("could not create watcher")
	}
	defer watcher.Close()
	defer producer.rabbitConn.Close()
	done, errs, files := make(chan bool), make(chan error, 1), make(chan Core.DataPacket, 1)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Create {
					file, err := parseFileFromEvent(event.Name)
					files <- file
					if err != nil {
						errs <- err
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	for _, element := range producer.config.Readers {
		err = watcher.Add(element.Directory)
		if err != nil {
			log.Fatal(err)
		}
	}
	go producer.insertToQueue(files)
	go handleExceptions(errs)                             // Handle File exceptions
	go handleExceptions(producer.rabbitConn.ErrorChannel) // Handle Rabbit exceptions
	<-done
}

func handleExceptions(errors chan error) {
	for err := range errors {
		log.Println(err)
	}
}

func (producer FileProducer) insertToQueue(packets chan Core.DataPacket) {
	for packet := range packets {
		if packet.Identifier != "" {
			producer.rabbitConn.SendMessage(producer.config.QueueName, packet)
		}
	}
}

func parseFileFromEvent(fileName string) (Core.DataPacket, error) {
	packet := Core.DataPacket{Identifier: ""}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return packet, err
	}
	rmErr := os.Remove(fileName)
	if rmErr != nil {
		return packet, rmErr
	}
	return Core.DataPacket{Identifier: fileName, Data: file}, nil
}
