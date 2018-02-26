package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	endpoint "github.com/CzarSimon/go-endpoint"
)

// Common infrastructure variables
var (
	configFile = filepath.Join(os.Getenv("HOME"), ".sws", "config.json")
)

// Auth authorization credentials for the sws-apiserver.
type Auth struct {
	AccessKey string `json:"accessKey"`
}

// Config for communicating with the sws-apiserver.
type Config struct {
	Auth Auth                `json:"auth"`
	API  endpoint.ServerAddr `json:"api"`
}

// GetConfig gets configuration.
func GetConfig() (Config, error) {
	var config Config
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Save stores configuration.
func (config Config) Save() error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, content, 0666)
}
