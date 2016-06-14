package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseConfig(configPath string) (*Config, error) {
	c := &Config{}
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, c)
	return c, err
}
