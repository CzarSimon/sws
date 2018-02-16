package service

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ReadService reads a service manifest from supplied file.
func ReadService(filename string) (Manifest, error) {
	var manifest Manifest
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return manifest, err
	}
	err = yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		return manifest, err
	}
	return manifest, nil
}
