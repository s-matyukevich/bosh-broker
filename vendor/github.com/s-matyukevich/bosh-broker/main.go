package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"

	"github.com/s-matyukevich/bosh-broker/source/broker"
	"github.com/s-matyukevich/bosh-broker/source/config"
)

func main() {
	logger := lager.NewLogger(appName(os.Getenv("VCAP_APPLICATION"), "bosh-broker"))
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	config, err := config.ParseConfig("config.yml")
	if err != nil {
		logger.Fatal("config", err)
		return
	}
	logger.Debug("config", lager.Data{"config": config})
	credentials := brokerapi.BrokerCredentials{
		Username: config.ServiceUser,
		Password: config.ServicePassword,
	}

	handler, err := broker.NewHandler(config)
	if err != nil {
		logger.Fatal("handler", err)
	}

	http.Handle("/", brokerapi.New(handler, logger, credentials))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Info("boot-up", lager.Data{"port": port})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("listen-and-serve", err)
	}
}

func appName(envJSON string, defaultName string) string {
	env := struct {
		ApplicationName string `json:"application_name"`
	}{}

	if envJSON == "" {
		goto DEFAULT
	}

	if err := json.Unmarshal([]byte(envJSON), &env); err != nil {
		panic(err)
	}

	if env.ApplicationName == "" {
		goto DEFAULT
	}

	return env.ApplicationName

DEFAULT:
	return defaultName
}
