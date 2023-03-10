package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Server config
	// Database config
	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`
	// Auth config
	Auth struct {
		JWT struct {
			// JWT secret key
			SecretKey string `yaml:"secret_key"`
			// JWT expiration time in seconds
			ExpirationTime int `yaml:"expiration_time"`
			Issuer         string
			Audience       string
			Subject        string
			JwtID          string
		} `yaml:"jwt"`
		SignInUrl string `yaml:"sign_in_url"`
	} `yaml:"auth"`
}

func LoadConfig(configYml []byte) *Config {
	return loadConfigYaml(configYml)
}

func loadConfigYaml(configYml []byte) *Config {
	// load yml from config.yml
	config := Config{}
	// unmarshal yml into config struct
	err := yaml.Unmarshal(configYml, &config) // TODO: this isn't working
	if err != nil {
		log.Fatalf("Error unmarshaling config.yml: %v", err)
	}
	return &config
}

func readConfigFile() []byte {
	// read config.yml
	configYml, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}
	return configYml
}
