package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type AppServer struct {
	Port           string `mapstructure:"Port"`
	ServiceName    string `mapstructure:"ServiceName"`
	ServiceVersion string `mapstructure:"ServiceVersion"`
}

type Database struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Name     string `mapstructure:"Name"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
}

type Cache struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Database int    `mapstructure:"Database"`
}

type Kafka struct {
	Broker       string `mapstructure:"Broker"`
	LoanTopic    string `mapstructure:"LoanTopic"`
	PaymentTopic string `mapstructure:"PaymentTopic"`
	Timeout      int    `mapstructure:"Timeout"`
}

type Config struct {
	AppServer AppServer `mapstructure:"AppServer"`
	Database  Database  `mapstructure:"Database"`
	Cache     Cache     `mapstructure:"Cache"`
	Kafka     Kafka     `mapstructure:"Kafka"`
}

func NewConfig(service string) (*Config, error) {
	var config Config

	configName := fmt.Sprintf("config-%s", service)
	viper.SetConfigName(configName)
	viper.AddConfigPath("./config-file")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
