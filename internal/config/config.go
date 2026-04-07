package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	cleanenvport "github.com/wb-go/wbf/config/cleanenv-port"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Storage StorageConfig `yaml:"storage"`
	Logger  LoggerConfig  `yaml:"logger"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"image-processor"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"local"`
}

type HTTPConfig struct {
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"image-tasks"`
	Group   string   `yaml:"group" env:"KAFKA_GROUP" env-default:"image-processor-group"`
}

type StorageConfig struct {
	UploadsDir   string `yaml:"uploads_dir" env:"STORAGE_UPLOADS_DIR" env-default:"./uploads"`
	ProcessedDir string `yaml:"processed_dir" env:"STORAGE_PROCESSED_DIR" env-default:"./processed"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

func Load() (*Config, error) {
	var cfg Config

	configPath := "config.yaml"

	err := cleanenvport.LoadPath(configPath, &cfg)
	if err != nil {
		fmt.Printf("Warning: config file not found or invalid: %v. Using defaults and env vars.\n", err)
		if envErr := cleanenv.ReadEnv(&cfg); envErr != nil {
			fmt.Printf("Error parsing env vars: %v\n", envErr)
		}
	}

	return &cfg, nil
}
