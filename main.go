package main

//
//import (
//	"log"
//
//	"github.com/YiftachE/DockerHPC/FileInputReader"
//	"github.com/YiftachE/DockerHPC/Middleware"
//	"github.com/YiftachE/DockerHPC/Core"
//
//)
//
//func main() {
//	done := make(chan bool)
//	prod, err := FileInputReader.CreateProducer("/home/yiftach/config.yml",Core.NewRedisWrapper())
//	if err != nil {
//		log.Fatal(err)
//	}
//	go prod.Start()
//	logging, err1 := Middleware.NewLoggingMiddleware("/home/yiftach/log.log")
//	logging2, err2 := Middleware.NewLoggingMiddleware("/home/yiftach/log2.log")
//	if err1 != nil || err2 != nil {
//		log.Fatalln(err1, err2)
//	}
//	middlewarePipeline := Middleware.NewPipeline("amqp://rabbitmq:rabbitmq@localhost:5672/",Core.NewRedisWrapper())
//	middlewarePipeline.Use(&logging)
//	middlewarePipeline.Use(&logging2)
//	go middlewarePipeline.Start()
//	<-done
//}

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/client"
)

func main() {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)
	// set "/foo" key with "bar" value
	log.Print("Setting '/foo' key with 'bar' value")
	resp, err := kapi.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
	// get "/foo" key's value
	log.Print("Getting '/foo' key value")
	resp, err = kapi.Get(context.Background(), "/foo", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
}
