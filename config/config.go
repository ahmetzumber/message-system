package config

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

type Config struct {
	AppName       string
	Port          int
	MongoDBConfig *MongoDBConfig
	RedisConfig   *RedisConfig
	WebhookConfig *WebhookConfig
}

type MongoDBConfig struct {
	URI        string
	Database   string
	Collection string
}

type RedisConfig struct {
	URI string
}

type WebhookConfig struct {
	BaseURL string
}

func readConfig(configPath, filename string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(filename)
	err := v.ReadInConfig()
	return v, err
}

func New(configPath, filename string) (*Config, error) {
	v, err := readConfig(configPath, filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	config := &Config{}
	return config, v.Unmarshal(config)
}

func (c *Config) Print() {
	pp.Println(c)
}
