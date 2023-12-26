package config

import (
	"io/ioutil"
)

const (
	DefaultDSN           = "dsn://"
	DefaultMigrationsDir = "file://migrations/auth"
)

type DbConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	Database DbConfig `yaml:"database"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		Database: DbConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
