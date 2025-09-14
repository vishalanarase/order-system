package config

import "github.com/spf13/viper"

type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

type AppConfig struct {
	HTTP        HTTPConfig `mapstructure:"http"`
	Application string     `mapstructure:"application"`
	Environment string     `mapstructure:"environment"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("http.port", "8080")
	viper.SetDefault("environment", "development")
	viper.SetDefault("application", "order-processing-system")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
