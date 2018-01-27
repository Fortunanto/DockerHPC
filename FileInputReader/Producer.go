package FileInputReader

import (
	"io/ioutil"
	"log"

	"github.com/YiftachE/DockerHPC/Core"
	"github.com/YiftachE/DockerHPC/Core/Interfaces"
	"github.com/fsnotify/fsnotify"
)

type FileProducer struct {
	config     FileInputReaderConfiguration
	rabbitConn Core.RabbitConnection
	shouldRun  bool
	files      chan Core.DataPacket
	errs       chan error
}

func (producer *FileProducer) Start() {
	defer producer.rabbitConn.Close()
	done := make(chan bool)
	producer.errs, producer.files = make(chan error, 100), make(chan Core.DataPacket, 100)
	defer producer.Stop()
	go producer.Insert(producer.files)
	go producer.HandleExceptions(producer.errs)
	go producer.HandleExceptions(producer.rabbitConn.ErrorChannel)
	for _, directory := range producer.config.Directories {
		go func(directory SourceConfiguration) {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				producer.errs <- err
				return
			}
			defer watcher.Close()
			watcher.Add(directory.Path)
			for producer.shouldRun {
				select {
				case event := <-watcher.Events:
					if event.Op == fsnotify.Create {
						results, err := parseFileFromEvent(event.Name, directory.Queues)
						for _, packet := range results {
							producer.files <- packet
						}
						if err != nil {
							producer.errs <- err
						}
					}
				case err := <-watcher.Errors:
					producer.errs <- err
				}
			}
		}(directory)
	}
	<-done
}

func (producer *FileProducer) HandleExceptions(errors chan error) {
	for err := range errors {
		log.Println(err)
	}
}

func (producer *FileProducer) Insert(packets chan Core.DataPacket) {
	for packet := range packets {
		if packet.Identifier != "" {
			packet.QueueName = packet.QueueName + "_input"
			producer.rabbitConn.SendMessage(Core.INSTREAMNAME, packet)
		}
	}
}

func CreateProducer(path string, store Interfaces.BackingStore) (FileProducer, error) {
	config, err := GetConfig(path)
	if err != nil {
		return FileProducer{}, err
	}
	conn, connectionError := Core.NewRabbitConnection(config.RabbitConnectionString, store)
	if connectionError != nil {
		return FileProducer{}, connectionError
	}
	return FileProducer{config: config, rabbitConn: conn, shouldRun: true}, nil
}

func (producer FileProducer) Stop() {
	producer.shouldRun = false
	close(producer.files)
}

func parseFileFromEvent(fileName string, queues []string) ([]Core.DataPacket, error) {
	log.Println("handling file :", fileName, queues)
	res := make([]Core.DataPacket, 0)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	for _, queue := range queues {
		res = append(res, Core.DataPacket{Identifier: fileName,
			QueueName: queue,
			Data: Core.NewPacketData(func() ([]byte, error) {
				return file, nil
			})})
	}
	return res, nil
}
