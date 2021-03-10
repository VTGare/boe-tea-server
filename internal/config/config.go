package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Config struct {
	Port     string `json:"port"`
	Database DB     `json:"database"`
}

type DB struct {
	URI       string `json:"uri"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	DefaultDB string `json:"default_db"`
}

func (c *Config) ConnectionURI() string {
	if c.Database.URI != "" {
		return c.Database.URI
	}

	var sb strings.Builder
	sb.WriteString("mongodb://")
	if c.Database.Username != "" {
		sb.WriteString(fmt.Sprintf("%v:%v@", c.Database.Username, c.Database.Password))
	}

	sb.WriteString(fmt.Sprintf("%v:%v", c.Database.Host, c.Database.Port))
	if c.Database.DefaultDB != "" {
		sb.WriteString("/" + c.Database.DefaultDB)
	}

	return sb.String()
}

//FromFile loads a configuration from a json file.
func FromFile(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

//FromEnv loads a configuration from environment variables, throwing an error when essential environment variables do not exist.
func FromEnv() (*Config, error) {
	//TODO: Create a ENV based configuration

	return &Config{}, nil
}
