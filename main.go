package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/didip/mcrouter-hub/application"
	"net/http"
	"os"
	"runtime"
)

func init() {
	logLevelString := os.Getenv("MCRHUB_LOG_LEVEL")
	if logLevelString == "" {
		logLevelString = "info"
	}
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err == nil {
		logrus.SetLevel(logLevel)
	}
}

// main runs the web server for resourced.
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	configFile := os.Getenv("MCROUTER_CONFIG_FILE")
	mcRouterHubCentralURL := os.Getenv("MCRHUB_CENTRAL_URL")

	if configFile == "" && mcRouterHubCentralURL != "" {
		err := errors.New("MCROUTER_CONFIG_FILE is required.")
		logrus.Fatal(err)
	}

	app, err := application.New(configFile, mcRouterHubCentralURL)
	if err != nil {
		logrus.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	app.WatchMcRounterConfigFile()
	app.ReportToCentral()

	httpAddr := os.Getenv("MCRHUB_ADDR")
	if httpAddr == "" {
		if app.IsCentral() {
			httpAddr = ":5002"
		} else {
			httpAddr = ":5001"
		}
	}

	httpsCertFile := os.Getenv("MCRHUB_CERT_FILE")
	httpsKeyFile := os.Getenv("MCRHUB_KEY_FILE")

	if httpsCertFile != "" && httpsKeyFile != "" {
		logrus.WithFields(logrus.Fields{
			"httpAddr": httpAddr,
		}).Info("Running HTTPS server")

		err = http.ListenAndServeTLS(httpAddr, httpsCertFile, httpsKeyFile, middle)
		if err != nil {
			logrus.Fatal(err)
		}

	} else {
		logrus.WithFields(logrus.Fields{
			"httpAddr": httpAddr,
		}).Info("Running HTTP server")

		err = http.ListenAndServe(httpAddr, middle)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
