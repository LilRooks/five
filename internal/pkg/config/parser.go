package config

import (
	"io/ioutil"

	"olympos.io/encoding/edn"
)

// ConfSet is a struct representing the configuration
type ConfSet struct {
	// IPFS
	RepoDir string
	RepoLDB string
	// Store
	Address string
	AddrDoc string
	MkStore bool
	MkLocal bool
	IDToken string
	// Commands
	SortsBy string
	AnonPut bool
}

// Read parses the configuration given at path.
func Read(path string) (*ConfSet, error) {
	if len(path) == 0 {
		return &ConfSet{}, nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &ConfSet{}, err
	}
	var ret *ConfSet
	err = edn.Unmarshal(data, &ret)
	return ret, nil
}

// Defaults parses the defaults if needed
func (c *ConfSet) Defaults(cArgs *ConfSet) error {
	// set overrides for each field
	if cArgs.MkStore {
		c.MkStore = true
	}
	if cArgs.MkLocal {
		c.MkLocal = true
	}
	return nil
}
