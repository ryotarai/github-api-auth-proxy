package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Passwords map[string][]string // username: [bcrypted password...]
}

func LoadYAMLFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.UnmarshalStrict(b, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
