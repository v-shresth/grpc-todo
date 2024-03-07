package utils

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type EnvConfig interface {
	GetJwtSecret() string
	GetServerPort() string
	GetMongoURI() string
	GetKafkaHost() string
	GetMailChimpApiKey() string
	GetSenderEmailAddress() string
}

type config struct {
	JwtSecret          string `env:"JWT_SECRET"`
	ServerPort         string `env:"PORT"`
	MongoUri           string `env:"MONGO_URI"`
	KafKaHost          string `env:"KAFKA_HOST"`
	MailChimpApiKey    string `env:"MAIL_CHIMP_API_KEY"`
	SenderEmailAddress string `env:"SENDER_EMAIL_ADDRESS"`
}

func NewEnvConfig() (EnvConfig, error) {
	envConfig := config{}
	if err := env.Parse(&envConfig); err != nil {
		return nil, fmt.Errorf("failed to load env config with error: %+v", err)
	}
	if envConfig.ServerPort == "" {
		envConfig.ServerPort = ":8080"
	}
	return &envConfig, nil
}

func (e *config) GetJwtSecret() string {
	if e == nil {
		return ""
	}
	return e.JwtSecret
}

func (e *config) GetServerPort() string {
	if e == nil {
		return ""
	}
	return e.ServerPort
}

func (e *config) GetMongoURI() string {
	if e == nil {
		return ""
	}
	return e.MongoUri
}

func (e *config) GetKafkaHost() string {
	if e == nil {
		return ""
	}
	return e.KafKaHost
}

func (e *config) GetMailChimpApiKey() string {
	if e == nil {
		return ""
	}
	return e.MailChimpApiKey
}

func (e *config) GetSenderEmailAddress() string {
	if e == nil {
		return ""
	}
	return e.SenderEmailAddress
}
