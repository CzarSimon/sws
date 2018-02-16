package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	yamlSuffixes = []string{".yml", ".yaml"}
	jsonSuffix   = ".json"
)

// ReadService reads a service manifest from supplied file.
func ReadService(filename string) (Manifest, error) {
	var manifest Manifest
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return manifest, err
	}
	err = deserializeService(filename, data, &manifest)
	if err != nil {
		return manifest, err
	}
	return manifest, nil
}

// deserializeService deserializes service definition if passed file is of supported type.
func deserializeService(filename string, data []byte, v interface{}) error {
	if isYaml(filename) {
		return yaml.Unmarshal(data, v)
	}
	if strings.HasSuffix(filename, jsonSuffix) {
		return json.Unmarshal(data, v)
	}
	return fmt.Errorf("Unkown filetype: %s", filename)
}

// isYaml checks if a given filename point to a yaml file.
func isYaml(filename string) bool {
	for _, yamlSuffix := range yamlSuffixes {
		if strings.HasSuffix(filename, yamlSuffix) {
			return true
		}
	}
	return false
}
