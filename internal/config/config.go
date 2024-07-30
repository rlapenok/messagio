package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/rlapenok/messagio/internal/utils"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}
type DataBaseConfig struct {
	Host     string `yaml:"host"`
	Db       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Url   []string `yaml:"url"`
	Topic string   `yaml:"topic"`
}

type Config struct {
	ServerConfig   ServerConfig   `yaml:"server"`
	DataBaseConfig DataBaseConfig `yaml:"database"`
	KafkaConfig    KafkaConfig    `yaml:"kafka"`
}

func New() *Config {

	var config Config
	var err error
	err = godotenv.Load("./.env")
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	path := os.Getenv("CONFIG_PATH")
	if path == "" {

		utils.Logger.Fatal("Insert path to config in .env")
	}
	err = cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
