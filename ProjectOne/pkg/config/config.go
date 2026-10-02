package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Jwt      JWTConfig      `mapstructure:"jwt"`
	Features FeaturesConfig `mapstructure:"features"`
	Quests   QuestsConfig   `mapstructure:"quests"`
}

type StorageConfig struct {
	PathSSL         string        `mapstructure:"path_ssl"`
	PathNotSSL      string        `mapstructure:"path_not_sll"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type ServerConfig struct {
	Port        int           `mapstructure:"port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type RedisConfig struct {
}

type FeaturesConfig struct {
	EnableSwagger bool `mapstructure:"enable_swagger"`
	EnableQuests  bool `mapstructure:"enable_quests"`
	EnableCache   bool `mapstructure:"enable_cache"`
	EnableMetrics bool `mapstructure:"enable_metrics"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	TTL    time.Duration `mapstructure:"ttl"`
}

type QuestsConfig struct {
}

func Load() (*Config, error) {
	const op = "config.Load"

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("C:/Users/akimo/Desktop/Something new/ProjectOne/pkg/config") //текущая директория

	setDefault()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read config file:%s,  %w", op, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config:%s,  %w", op, err)
	}

	return &cfg, nil
}

func setDefault() {
	// Server defaults
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.timeout", 10*time.Second)
	viper.SetDefault("server.idle_timeout", 30*time.Second)

	// Redis defaults
	// viper.SetDefault("redis.url", "redis://localhost:6379")
	// viper.SetDefault("redis.password", "")
	// viper.SetDefault("redis.db", 0)
	// viper.SetDefault("redis.timeout", 3*time.Second)

	// JWT defaults
	viper.SetDefault("jwt.secret", "super-secret-jwt-key")
	viper.SetDefault("jwt.ttl", 24*time.Hour)

	// Features defaults
	viper.SetDefault("features.enable_swagger", true)
	viper.SetDefault("features.enable_quests", true)
	viper.SetDefault("features.enable_cache", true)
	viper.SetDefault("features.enable_metrics", false)
}
