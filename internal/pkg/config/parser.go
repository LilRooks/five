package config

import (
	"io/ioutil"

	"olympos.io/encoding/edn"
)

// ConfSet is a struct representing the configuration
type ConfSet struct {
	Address string
	KeyPath string
}

// Read parses the configuration given at path.
func Read(path string) (ConfSet, error) {
	if len(path) == 0 {
		return ConfSet{}, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ConfSet{}, err
	}
	var ret ConfSet
	err = edn.Unmarshal(data, &ret)
	return ret, nil
}
