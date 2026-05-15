package core_postgres_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required: "true"`
	Port     string        `envconfig:"PORT" default: "5432"`
	User     string        `envconfig:"USER" required: "true"`
	Password string        `envconfig:"PASSWORD" required: "true"`
	Database string        `envconfig:"DB" required: "true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required: "true"`
}

func NewCfg() (Config, error) {
	var cfg Config

	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewCfgMust() Config {
	cfg, err := NewCfg()

	if err != nil {
		err = fmt.Errorf("get postgres connection pool config: %w", err)
		panic(err)
	}

	return cfg
}
