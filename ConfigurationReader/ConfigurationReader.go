package ConfigurationReader

import (
	"encoding/json"
	"os"
)

type ConfigurationReader struct{}

type Configuration struct {
	clientKey    string
	clientSecret string
}

var configuration = Configuration{}

func (c *ConfigurationReader) readConfiguration(path string) (Configuration, error) {

	file, err := os.Open("config.json")
	if err != nil {
		return Configuration{}, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
