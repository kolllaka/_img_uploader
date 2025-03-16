package model

import "time"

type Config struct {
	DB struct {
		Path string `yaml:"path"`
	} `yaml:"db"`
	Files struct {
		ClearDelay time.Duration `yaml:"clear_delay"`
	} `yaml:"files"`
	Images struct {
		SaveFolder   string        `yaml:"save_folder"`
		DelayExpires time.Duration `yaml:"delay_expires"`
		Width        uint          `yaml:"width"`
		Height       uint          `yaml:"height"`
	} `yaml:"images"`
}

func NewConfig() *Config {
	return &Config{}
}
