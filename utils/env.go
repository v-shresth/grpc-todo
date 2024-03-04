package utils

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	jwtSecret     string `env:"JWT_SECRET"`
	serverPort    string `env:"PORT"`
	mongoUri      string `env:"MONGO_URI"`
	kafKaProvider string `env:"KAFKA_PROVIDER"`
}

func NewEnvConfig() (*EnvConfig, error) {
	config := EnvConfig{}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to load env config with error: %+v", err)
	}
	return &config, nil
}

func (e *EnvConfig) GetJwtSecret() string {
	return e.jwtSecret
}

func (e *EnvConfig) GetServerPort() string {
	if e.serverPort == "" {
		return ":8080"
	}
	return e.serverPort
}

func (e *EnvConfig) GetMongoURI() string {
	return e.mongoUri
}

func (e *EnvConfig) GetKafkaHost() string {
	return e.kafKaProvider
}
