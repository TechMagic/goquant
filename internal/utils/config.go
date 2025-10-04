package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config is reserved for future extensions (currently a simple example)
type Config struct {
	DataDir string `yaml:"data_dir"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}
