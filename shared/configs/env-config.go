package configs

import (
	"log"

	"github.com/spf13/viper"
)

type EnvConfigs struct {
	Port            string `mapstructure:"PORT"`
	DbPassword      string `mapstructure:"DB_PASSWORD"`
	DatabaseUrl     string `mapstructure:"DATABASE_URL"`
	DebankAccessKey string `mapstructure:"DEBANK_ACCESS_KEY"`
}

var EnvConfigVars *EnvConfigs

func InitEnvConfigs() {
	var err error
	EnvConfigVars, err = loadEnvVariables()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
}

func loadEnvVariables() (config *EnvConfigs, err error) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env variables into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return config, nil
}

// GetDebankAccessKey returns the value of DEBANK_ACCESS_KEY
func (env *EnvConfigs) GetDebankAccessKeyHeader() string {
	return env.DebankAccessKey
}
