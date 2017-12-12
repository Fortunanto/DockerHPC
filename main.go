package main

import (
	"github.com/YiftachE/DockerHPC/SourceProducer"
	"log"
)

func main() {
	r, err := SourceProducer.CreateProducer("./SourceProducer/config.yml")
	if err != nil {
		log.Fatal(err)
	}else{
		r.Start()
	}
}
