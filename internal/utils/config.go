package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config 为今后扩展保留（当前简单示例）
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
