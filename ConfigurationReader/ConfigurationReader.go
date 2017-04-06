package ConfigurationReader

import (
	"encoding/hex"
	"encoding/json"
	"os"
)

type Configuration struct {
	SecureCookieSecretHex string `json:"secureCookieSecretHex"`

	ClientKey    string `json:"clientKey"`
	ClientSecret string `json:"clientSecret"`
}

var (
	configuration     Configuration
	configurationRead = false
)

func GetConfiguration() (*Configuration, error) {
	if !configurationRead {
		if err := ReadConfiguration("config.json", &configuration); err != nil {
			return nil, err
		}
		configurationRead = true
	}
	return &configuration, nil
}

func ReadConfiguration(path string, config *Configuration) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}

	return nil
}

func (c Configuration) GetAuthCookieSecret() ([]byte, error) {
	return hex.DecodeString(c.SecureCookieSecretHex)
}
