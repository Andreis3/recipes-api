package utils

import "github.com/spf13/viper"

type Config struct {
	Port            string `mapstructure:"PORT"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	MongoDB         string `mapstructure:"MONGO_DB"`
	MongoCollection string `mapstructure:"MONGO_COLLECTION"`
	RedisURI        string `mapstructure:"REDIS_URI"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	RedisDB         int    `mapstructure:"REDIS_DB"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
