package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string       `mapstructure:"env"`
	Server      ServerConfig `mapstructure:"server"`
	StoragePath string       `mapstructure:"storage_path"`
	Redis       RedisConfig  `mapstructure:"redis"`
	Jwt         string       `mapstructure:"jwt"`
	Features    RedisConfig  `mapstructure:"features"`
	Quests      QuestsConfig `mapstructure:"quests"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Timeout      time.Duration `mapstructure:"timeout"`
	IddleTimeout time.Duration `mapstructure:"iddle_timeout"`
}

type RedisConfig struct {
}

type FeaturesConfig struct {
	EnableSwagger bool `mapstructure:"enable_swagger"`
	EnableQuests  bool `mapstructure:"enable_quests"`
	EnableCache   bool `mapstructure:"enable_cache"`
	EnableMetrics bool `mapstructure:"enable_metrics"`
}

type QuestsConfig struct {
}

func Load() *Config {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("C:/Users/akimo/Desktop/Something new/ProjectOne/pkg/config") //текущая директория

	setDefault()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Warning: Config file not found, %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}

func setDefault() {
	// Server defaults
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("server.read_timeout", 10*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)

	// Redis defaults
	// viper.SetDefault("redis.url", "redis://localhost:6379")
	// viper.SetDefault("redis.password", "")
	// viper.SetDefault("redis.db", 0)
	// viper.SetDefault("redis.timeout", 3*time.Second)

	// JWT defaults
	// viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
	// viper.SetDefault("jwt.expires_in", 24*time.Hour)

	// Features defaults
	viper.SetDefault("features.enable_swagger", true)
	viper.SetDefault("features.enable_quests", true)
	viper.SetDefault("features.enable_cache", true)
	viper.SetDefault("features.enable_metrics", false)
}
