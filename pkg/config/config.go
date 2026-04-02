package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// MustLoad reads a YAML file into out and panics on error.
// This keeps the scaffold code compact for early-stage setup.
func MustLoad(path string, out any) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(b, out); err != nil {
		panic(err)
	}
}
