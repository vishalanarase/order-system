package config

import "github.com/spf13/viper"

type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	GroupID  string   `mapstructure:"group_id"`
	Topics   []string `mapstructure:"topics"`
	ClientID string   `mapstructure:"client_id"`
}

type AppConfig struct {
	Kafka       KafkaConfig `mapstructure:"kafka"`
	HTTP        HTTPConfig  `mapstructure:"http"`
	Application string      `mapstructure:"application"`
	Environment string      `mapstructure:"environment"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("http.port", "8080")
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.group_id", "order-service-group")
	viper.SetDefault("kafka.topics", []string{"orders", "payments", "notifications"})
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
