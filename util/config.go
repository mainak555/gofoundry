package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Configure loads environment-specific configuration into T using Viper.
func Configure[T interface{}](configPath string) (settings *T, err error) {
	envName := os.Getenv("NODE_ENV")
	if envName == "" {
		envName = "local"
		os.Setenv("NODE_ENV", "local")
	}

	if godotenv.Load(os.Getenv("NODE_ENV")+".env") != nil {
		fmt.Println(fmt.Sprintf("no %v.env", os.Getenv("NODE_ENV")))
	}

	viper.AddConfigPath(configPath + "/")
	viper.SetConfigType("yml")
	viper.SetConfigName(envName)
	viper.AutomaticEnv()

	/*var defaultConfiguration []byte
	.ReadConfig(bytes.NewBuffer(defaultConfiguration))*/

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprint("Reading File Error!", err))
	} else if err := viper.Unmarshal(&settings); err != nil {
		return nil, errors.New(fmt.Sprint("Fatal Error!", err))
	}
	return
}
