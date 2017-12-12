package SourceProducer

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type FileConfig struct {
	Directory   string
	PollingRate int
}
type FileProducerConfig struct {
	RabbitConnectionString string
	QueueName              string
	Readers                []FileConfig
}

func GetConfig(configPath string) (FileProducerConfig, error) {
	prodConfig := FileProducerConfig{}
	file, fileErr := ioutil.ReadFile(configPath)
	if fileErr != nil {
		return prodConfig, fileErr
	}
	err := yaml.Unmarshal(file, &prodConfig)
	if err != nil {
		return prodConfig, err
	}
	return prodConfig, nil
}
