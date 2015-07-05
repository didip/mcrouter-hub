package main

import (
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
	app, err := application.New()
	if err != nil {
		logrus.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	app.CollectData()
	app.ReportConfigToCentral()
	app.ReportStatsToCentral()
	app.ReportStatsToNewrelicInsights()

	httpAddr := app.Settings["MCRHUB_ADDR"]
	httpsCertFile := app.Settings["MCRHUB_CERT_FILE"]
	httpsKeyFile := app.Settings["MCRHUB_KEY_FILE"]

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
