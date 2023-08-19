package config

import (
	"log"

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
		// Kratos config
		Kratos struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		} `yaml:"kratos"`
		SignInUrl string `yaml:"sign_in_url"`
	} `yaml:"auth"`
}

func LoadConfig(configYml []byte) *map[string]Config {
	return loadConfigYaml(configYml)
}

func loadConfigYaml(configYml []byte) *map[string]Config {
	// load yml from config.yml
	config := map[string]Config{}
	// unmarshal yml into config struct
	err := yaml.Unmarshal(configYml, &config) // TODO: this isn't working
	if err != nil {
		log.Fatalf("error unmarshaling config.yml: %v", err)
	}
	return &config
}
