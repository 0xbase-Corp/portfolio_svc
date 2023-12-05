package configs

import (
	"log"

	"github.com/spf13/viper"
)

var envConfigVars *EnvConfigs

func InitEnvConfigs() {
	var err error
	envConfigVars, err = loadEnvVariables()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
}
func GetDatabaseUrl() string {
	return envConfigVars.DatabaseUrl
}
func GetPort() string {
	return envConfigVars.Port
}

type EnvConfigs struct {
	Port        string `mapstructure:"PORT"`
	DbPassword  string `mapstructure:"DB_PASSWORD"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

// rest of the code remains the same

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
	return
}
