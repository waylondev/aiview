// Package config provides configuration loading with viper, supporting YAML files and environment variables.
package config

import (
	"github.com/spf13/viper"
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
type PlatformsConfig struct {
	Bilibili    PlatformConfig `mapstructure:"bilibili"`
	Douyin      PlatformConfig `mapstructure:"douyin"`
	Weibo       PlatformConfig `mapstructure:"weibo"`
	Kuaishou    PlatformConfig `mapstructure:"kuaishou"`
	Xiaohongshu PlatformConfig `mapstructure:"xiaohongshu"`
	Zhihu       PlatformConfig `mapstructure:"zhihu"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Platform: "bilibili",
		CacheTTL: 300,
		Output:   "auto",
		Platforms: PlatformsConfig{
			Bilibili: PlatformConfig{
				Timeout: 30,
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
	viper.SetDefault("cache_ttl", 300)
	viper.SetDefault("output", "auto")
	viper.SetDefault("platforms.bilibili.timeout", 30)
	viper.SetDefault("platforms.douyin.timeout", 30)
	viper.SetDefault("platforms.weibo.timeout", 30)
	viper.SetDefault("platforms.kuaishou.timeout", 30)
	viper.SetDefault("platforms.xiaohongshu.timeout", 30)
	viper.SetDefault("platforms.zhihu.timeout", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, use defaults
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}