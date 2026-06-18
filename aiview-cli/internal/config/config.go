// Package config provides configuration loading with viper, supporting YAML files and environment variables.
package config

import (
	"errors"

	"github.com/spf13/viper"
)

const (
	// DefaultTimeout is the default HTTP request timeout in seconds.
	DefaultTimeout = 30
	// DefaultCacheTTL is the default cache time-to-live in seconds.
	DefaultCacheTTL = 300
)

// Config holds the global configuration for aiview.
type Config struct {
	Platform  string           `mapstructure:"platform"`
	CacheTTL  int              `mapstructure:"cache_ttl"`
	Output    string           `mapstructure:"output"`
	Platforms PlatformsConfig  `mapstructure:"platforms"`
}

// PlatformConfig holds configuration for a social media platform.
type PlatformConfig struct {
	Cookies string `mapstructure:"cookies"`
	Timeout int    `mapstructure:"timeout"`
}

// PlatformsConfig holds platform-specific configurations.
type PlatformsConfig map[string]PlatformConfig

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Platform: "bilibili",
		CacheTTL: DefaultCacheTTL,
		Output:   "auto",
		Platforms: PlatformsConfig{
			"bilibili": PlatformConfig{
				Timeout: DefaultTimeout,
			},
		},
	}
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.aiview")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("AIVIEW")
	viper.AutomaticEnv()

	// Bind platform-specific config defaults
	viper.SetDefault("platform", "bilibili")
	viper.SetDefault("cache_ttl", DefaultCacheTTL)
	viper.SetDefault("output", "auto")
	viper.SetDefault("platforms.bilibili.timeout", DefaultTimeout)
	viper.SetDefault("platforms.douyin.timeout", DefaultTimeout)
	viper.SetDefault("platforms.weibo.timeout", DefaultTimeout)
	viper.SetDefault("platforms.kuaishou.timeout", DefaultTimeout)
	viper.SetDefault("platforms.xiaohongshu.timeout", DefaultTimeout)
	viper.SetDefault("platforms.zhihu.timeout", DefaultTimeout)

	if err := viper.ReadInConfig(); err != nil {
		var configErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configErr) {
			return nil, err
		}
		// Config file not found, use defaults
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}