package initializers

import (
	"core/configuration"
	"core/log"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func InitializeConfiguration() configuration.Configuration {
	logger := log.NewLogger()
	var config configuration.Configuration

	viper.SetConfigType("json")
	viper.SetConfigFile("./configuration/configuration.json")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic(err)
	}

	err = viper.UnmarshalExact(&config)
	if err != nil {
		logger.Panic(err)
	}

	err = validator.New().Struct(config)
	if err != nil {
		logger.Panic(err)
	}

	return config
}
