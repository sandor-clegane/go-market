package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultServerEndpoint = "localhost:8080"
	DefaultAccrualSystemAddress
	DefaultKey = "SuperSecretKey2022"
	DefaultDatabaseDSN
)

type Config struct {
	ServerEndpoint       string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseAddress      string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Key                  string
}

func (c *Config) parseArgsCMD() {
	if !flag.Parsed() {
		flag.StringVar(&c.ServerEndpoint, "a",
			DefaultServerEndpoint, "http server launching address")
		flag.StringVar(&c.DatabaseAddress, "d", DefaultDatabaseDSN,
			"DB connection address")
		flag.StringVar(&c.AccrualSystemAddress, "r", DefaultAccrualSystemAddress,
			"address of accrual service")
		flag.Parse()
	}
}

func (c *Config) applyConfig(other Config) {
	if c.ServerEndpoint == DefaultServerEndpoint {
		c.ServerEndpoint = other.ServerEndpoint
	}
	if c.Key == DefaultKey {
		c.Key = other.Key
	}
	if c.DatabaseAddress == DefaultDatabaseDSN {
		c.DatabaseAddress = other.DatabaseAddress
	}
	if c.AccrualSystemAddress == DefaultAccrualSystemAddress {
		c.AccrualSystemAddress = other.AccrualSystemAddress
	}
}

func (c *Config) Init() error {
	var c2 Config
	//parsing env config
	err := env.Parse(c)
	if err != nil {
		return err
	}
	//parsing command line config
	c2.parseArgsCMD()
	//applying config
	c.applyConfig(c2)
	return nil
}
