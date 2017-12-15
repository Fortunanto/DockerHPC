package FileInputReader

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SourceConfiguration struct {
	Path   string
	PollingRate int    `yaml:",flow"`
	Queues   []string `yaml:queues`
}
type FileInputReaderConfiguration struct {
	RabbitConnectionString string `yaml:"rabbitConnectionString"`
	Directories                []SourceConfiguration
}

func GetConfig(configPath string) (FileInputReaderConfiguration, error) {
	prodConfig := FileInputReaderConfiguration{}
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
