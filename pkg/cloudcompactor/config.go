package cloudcompactor

import (
	"os"

	"github.com/lerenn/cloud-compactor/pkg/accessors/ftp"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Path  string
	Speed string

	Formats struct {
		Inputs          []string
		ProcessedSuffix string `yaml:"processed_suffix"`
		Output          string
	}

	FTP ftp.Config
}

func LoadConfigFromFile(path string) (Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
