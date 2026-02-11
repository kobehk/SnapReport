package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	DDPai struct {
		BaseURL        string `yaml:"base_url"`
		TimeoutSeconds int    `yaml:"timeout_seconds"`
		MockMode       bool   `yaml:"mock_mode"`
	} `yaml:"ddpai"`
	Geocoder struct {
		Type      string `yaml:"type"`       // "nominatim" 或 "amap"
		UserAgent string `yaml:"user_agent"` // 仅 Nominatim 使用
		APIKey    string `yaml:"api_key"`    // 仅 AMap 使用
	} `yaml:"geocoder"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	// Defaults
	cfg.Server.Port = 8080
	cfg.DDPai.BaseURL = "http://193.168.0.1"
	cfg.DDPai.TimeoutSeconds = 5
	cfg.Geocoder.Type = "nominatim"
	cfg.Geocoder.UserAgent = "SnapReport/1.0"
	cfg.Geocoder.APIKey = ""

	f, err := os.Open(path)
	if err != nil {
		log.Printf("Warning: config file %s not found, using defaults: %v", path, err)
		return &cfg, nil
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
