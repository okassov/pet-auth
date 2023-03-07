package config

import "github.com/spf13/viper"

type Config struct {
	PG     PGConfig     `mapstructure:",squash"`
	JWT    JWTConfig    `mapstructure:",squash"`
	Server ServerConfig `mapstructure:",squash"`
}

type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
}

type PGConfig struct {
	PGUrl      string `mapstructure:"PG_URL"`
	PGPort     string `mapstructure:"PG_PORT"`
	PGUser     string `mapstructure:"PG_USER"`
	PGPassword string `mapstructure:"PG_PASSWORD"`
	PGDatabase string `mapstructure:"PG_DATABASE"`
}

type JWTConfig struct {
	JWTKey string `mapstructure:"JWT_KEY"`
	JWTTtl int    `mapstructure:"JWT_TTL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
