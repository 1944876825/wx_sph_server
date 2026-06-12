package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Conf = Model{}

type Model struct {
	Port  int    `yaml:"port"`
	Mysql string `yaml:"mysql"`
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
	AI    AIConfig `yaml:"ai"`
}

type AIConfig struct {
	Enabled       bool   `yaml:"enabled"`
	APIKey        string `yaml:"api_key"`
	BaseURL       string `yaml:"base_url"`
	Model         string `yaml:"model"`
	DefaultPrompt string `yaml:"default_prompt"`
}

func Load() {
	content, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(content, &Conf); err != nil {
		panic(err)
	}
}
func Save() error {
	y, err := yaml.Marshal(Conf)
	if err != nil {
		return err
	}
	if err := os.WriteFile("./config.yaml", y, 0644); err != nil {
		return err
	}
	return nil
}
