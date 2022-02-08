package config

import "github.com/spf13/viper"

func App() map[string]interface{} {
	return viper.GetStringMap("app")
}

func AppName() string {
	return viper.GetString("app.name")
}

func AppMySQL() map[string]interface{} {
	return viper.GetStringMap("app.mysql")
}

func AppMySQLHost() string {
	return viper.GetString("app.mysql.host")
}

func AppMySQLPort() int {
	return viper.GetInt("app.mysql.port")
}

func AppMySQLUsername() string {
	return viper.GetString("app.mysql.username")
}

func AppMySQLPassword() string {
	return viper.GetString("app.mysql.password")
}

func AppMySQLDatabase() string {
	return viper.GetString("app.mysql.database")
}
