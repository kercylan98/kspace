package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	//language=yaml
	var configTemplate = []byte(`
app:
  name: "Unknown Service"
  mysql:
    host: 127.0.0.1
    port: 3306
    username: root
    password: root
    database: xxx
`)
	var initialized = false
initConfig:
	{
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil && !initialized {
			initialized = false
			if err = viper.ReadConfig(bytes.NewBuffer(configTemplate)); err != nil {
				panic(err)
			}
			if err = viper.WriteConfigAs("./config/config"); err != nil {
				panic(err)
			}
			goto initConfig
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}

	}

}
