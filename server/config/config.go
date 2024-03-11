package config

import "os"

type Config struct {
	Listen string
	Port   string
	Db     *Db
	Key    string
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	c.Listen = os.Getenv("SERVER_LISTEN")
	c.Port = os.Getenv("SERVER_PORT")
	c.Db = &Db{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}
	c.Key = os.Getenv("KEY")
}
