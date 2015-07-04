package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/didip/mcrouter-hub/application"
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
	mcRouterAddr := os.Getenv("MCROUTER_ADDR")
	mcRouterConfigFile := os.Getenv("MCROUTER_CONFIG_FILE")
	mcRouterHubCentralURL := os.Getenv("MCRHUB_CENTRAL_URL")
	mcRouterHubReadOnlyString := os.Getenv("MCRHUB_READ_ONLY")

	mcRouterHubReadOnly := true
	if mcRouterHubReadOnlyString == "false" {
		mcRouterHubReadOnly = false
	}

	if mcRouterConfigFile == "" && mcRouterHubCentralURL != "" {
		err := errors.New("MCROUTER_CONFIG_FILE is required.")
		logrus.Fatal(err)
	}

	nrInsightsURL := os.Getenv("NR_INSIGHTS_URL")
	if nrInsightsURL == "" {
		nrInsightsURL = "https://insights-collector.newrelic.com/v1/accounts/1/events"
	}
	nrInsightsInsertKey := os.Getenv("NR_INSIGHTS_INSERT_KEY")

	app, err := application.New(mcRouterHubReadOnly, mcRouterConfigFile, mcRouterAddr, mcRouterHubCentralURL, nrInsightsURL, nrInsightsInsertKey)
	if err != nil {
		logrus.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	app.CollectData()
	app.ReportToCentral()
	app.ReportToNewrelicInsights()

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
