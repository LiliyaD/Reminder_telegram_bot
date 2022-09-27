//go:build integration
// +build integration

package config

import (
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "QA"

type Config struct {
	Host                   string `default:":8081"`
	DbPort                 string `default:"6432"`
	DbHost                 string `default:"localhost"`
	User, Password, DBname string
}

func FromEnv() *Config {
	cfg := &Config{}

	var exists bool

	cfg.User, exists = os.LookupEnv("POSTGRESQL_USER")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_USER doesn't exist"))
	}

	cfg.Password, exists = os.LookupEnv("POSTGRESQL_PASS")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_PASS doesn't exist"))
	}

	cfg.DBname, exists = os.LookupEnv("POSTGRESQL_DB")
	if !exists {
		log.Fatal(errors.New("Environment variable POSTGRESQL_DB doesn't exist"))
	}

	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
