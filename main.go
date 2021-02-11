package main

import (
	"github.com/llamadeus/iot-logger/graph"
	"github.com/llamadeus/iot-logger/graph/generated"
	"github.com/llamadeus/iot-logger/internal/server"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	viper.SetDefault("APP_ENV", "production")
	viper.SetDefault("APP_PORT", "8080")
}

func main() {
	app := server.Init(server.Options{
		Port:          viper.GetString("APP_PORT"),
		Schema:        generated.NewExecutableSchema(graph.New()),
		AntiCSRF:      !isDevelopment(),
		Playground:    isDevelopment(),
		Introspection: isDevelopment(),
	})
	err := app.Start()

	if err != nil {
		log.Fatal(err)
	}
}

func isDevelopment() bool {
	return viper.GetString("APP_ENV") == "development"
}
