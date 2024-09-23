package config

import (
	"fmt"
	"io/ioutil"
	"youke/global/cos"
	"youke/global/database"
	"youke/global/logger"
	"youke/global/ocr"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql  *database.Mysql      `yaml:"Mysql"`
	Logger *logger.LoggerConfig `yaml:"Logger"`
	Cos    *cos.CosConfig       `yaml:"Cos"`
	Ocr    *ocr.OcrConfig       `yaml:"Ocr"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
