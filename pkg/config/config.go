package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	ListenAddr    string              `yaml:"listenAddr"`
	OriginURL     string              `yaml:"originURL"`
	OPAPolicyFile string              `yaml:"opaPolicyFile"`
	AccessToken   string              `yaml:"accessToken"`
	Passwords     map[string][]string `yaml:"passwords"` // username: [bcrypted password...]
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

	config.LoadFromEnv()

	return config, nil
}

func (c *Config) LoadFromEnv() {
	if accessToken := os.Getenv("GHPROXY_ACCESS_TOKEN"); accessToken != "" {
		c.AccessToken = accessToken
	}
}
