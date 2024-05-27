package configs

import (
	"embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed application.yaml
var f embed.FS
var config Config

type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Repository  string `yaml:"repository"`
}

func GetConfig() *Config {
	if config != (Config{}) {
		return &config
	}

	yamlFile, err := f.ReadFile("application.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalln(err)
	}

	return &config
}
