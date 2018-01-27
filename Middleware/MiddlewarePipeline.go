package Middleware

import (
	"log"

	"github.com/YiftachE/DockerHPC/Core"
	"github.com/YiftachE/DockerHPC/Core/Interfaces"
)

type Pipeline struct {
	rabbitConn  Core.RabbitConnection
	middlewares []Middleware
}

func NewPipeline(rabbitConnString string, store Interfaces.BackingStore) Pipeline {
	rabbitConn, rabbitErr := Core.NewRabbitConnection(rabbitConnString, store)
	
	if rabbitErr != nil {
		log.Fatal(rabbitErr)
	}

	return Pipeline{rabbitConn: rabbitConn}

}

func (pipeline *Pipeline) Use(middleware Middleware) {
	pipeline.middlewares = append(pipeline.middlewares, middleware)
}

func (pipeline *Pipeline) Start() {
	go pipeline.HandleExceptions(pipeline.rabbitConn.ErrorChannel)

	packets, err := pipeline.rabbitConn.ReceiveMessage(Core.INSTREAMNAME)
	if err != nil {
		log.Fatalln(err)
	}
	done := make(chan bool)
	go func() {
		for packet := range packets {
			for _, ware := range pipeline.middlewares {
				packet = ware.HandleMessage(packet)
			}
			pipeline.rabbitConn.SendMessage(packet.QueueName, packet)
		}
	}()
	<-done
}
func (pipeline *Pipeline) HandleExceptions(errors chan error) {
	for err := range errors {
		log.Println(err)
	}
}
