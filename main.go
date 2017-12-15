package main

import (
	"log"

	"github.com/YiftachE/DockerHPC/FileInputReader"
)

func main() {
	r, err := FileInputReader.CreateProducer("./FileInputReader/config.yml")
	if err != nil {
		log.Fatal(err)
	} else {
		r.Start()
	}
}
