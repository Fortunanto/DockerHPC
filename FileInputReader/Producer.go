package FileInputReader

import (
	"io/ioutil"
	"log"

	"github.com/YiftachE/DockerHPC/Core"
	"github.com/fsnotify/fsnotify"
)

type FileProducer struct {
	config     FileInputReaderConfiguration
	rabbitConn Core.RabbitConnection
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
	defer producer.rabbitConn.Close()
	done, errs, files := make(chan bool), make(chan error, 100), make(chan Core.DataPacket, 100)
	go producer.insertToQueue(files)
	go handleExceptions(errs)                             // Handle File exceptions
	go handleExceptions(producer.rabbitConn.ErrorChannel) // Handle Rabbit exceptions
	for _, directory := range producer.config.Directories {
		go func(directory SourceConfiguration) {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				errs <- err
				return
			}
			defer watcher.Close()
			watcher.Add(directory.Path)
			for {
				select {
				case event := <-watcher.Events:
					if event.Op == fsnotify.Create {
						results, err := parseFileFromEvent(event.Name, directory.Queues)
						for _, packet := range results {
							files <- packet
						}
						if err != nil {
							errs <- err
						}
					}
				case err := <-watcher.Errors:
					errs <- err
				}
			}
		}(directory)
	}
	<-done
}

func parseFileFromEvent(fileName string, queues []string) ([]Core.DataPacket, error) {
	log.Println("handling file :", fileName, queues)
	res := make([]Core.DataPacket, 0)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	for _, queue := range queues {
		res = append(res, Core.DataPacket{Identifier: fileName, QueueName: queue, Data: file})
	}
	return res, nil
}

func handleExceptions(errors chan error) {
	for err := range errors {
		log.Println(err)
	}
}

func (producer FileProducer) insertToQueue(packets chan Core.DataPacket) {
	for packet := range packets {
		if packet.Identifier != "" {
			producer.rabbitConn.SendMessage(packet.QueueName, packet)
		}
	}
}
