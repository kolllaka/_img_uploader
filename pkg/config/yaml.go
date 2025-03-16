package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadByPath(path string, cfg any) error {
	conFile, err := os.Open(path)
	if err != nil {

		return fmt.Errorf("could not open config file: %w", err)
	}
	defer conFile.Close()

	return yaml.NewDecoder(conFile).Decode(cfg)
}

func MustLoadByPath(path string, cfg any) {
	if err := LoadByPath(path, cfg); err != nil {
		panic(err)
	}
}
