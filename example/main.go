package main

import (
	"configs"
	"errors"
	"github.com/caarlos0/env/v11"
	"log"
	"os"
)

// EXAMPLE

type Config struct {
	Name     string `env:"MY_NAME"`
	Username string `env:"MY_USERNAME" secret:""`
	Password string `env:"MY_PASSWORD" secret:""`
}

var config Config

func (Config) Register() error {
	return env.Parse(&config)
}

func (Config) Validation() error {
	if config.Username == "" {
		return errors.New("username cannot be empty")
	}
	if config.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func (Config) Print() interface{} {
	return config
}

func main() {
	_ = os.Setenv("MY_NAME", "My First Configuration")
	_ = os.Setenv("MY_USERNAME", "testUserName")
	_ = os.Setenv("MY_PASSWORD", "testUserPassword")

	err := configs.Load(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("name", config.Name)
	log.Println("username: ", config.Username)
	log.Println("password: ", config.Password)

	log.Println("configs successfully loaded")
}
