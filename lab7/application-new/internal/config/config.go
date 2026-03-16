package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPathEnvKey = "CONFIG_PATH"
)

// Config represents the configuration structure
type Config struct {
	// Add your configuration fields here
	// Example:
	CacheConfig
	HTTPServer  `yaml:"http_server"`
}

type CacheConfig struct {
	URL string
}

type HTTPServer struct {
	Address     string        `yaml:"address"			env-default:"localhost:8080`
	Timeout     time.Duration `yaml:"timeout" 			env-required:"true"`
	IdleTimeout time.Duration `yaml:"iddle_timeout" 	env-required:"true"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	if err := secretViaEnv(&cfg); err != nil {
		log.Fatalf("failed to load secret: %s", err.Error())
	}

	return &cfg
}

func secretViaEnv(cfg *Config) error {
	cacheURL := os.Getenv("CACHE_URL")
	if cacheURL == "" {
		return os.ErrNotExist
	}
	cfg.CacheConfig.URL = cacheURL
	
	log.Print("secret has been loaded")
	return nil
}