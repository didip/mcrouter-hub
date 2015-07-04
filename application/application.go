package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/didip/mcrouter-hub/handlers"
	"github.com/didip/mcrouter-hub/middlewares"
	"github.com/didip/mcrouter-hub/models"
	"github.com/didip/mcrouter-hub/storage"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
)

func New(readOnly bool, mcRouterConfigFile, mcRouterAddr, mcRouterHubCentralURL, nrInsightsURL, nrInsightsInsertKey string) (*Application, error) {
	app := &Application{}
	app.ReadOnly = readOnly
	app.McRouterAddr = mcRouterAddr
	app.McRouterStatsManager = models.NewMcRouterStatsManager(mcRouterAddr)
	app.McRouterHubCentralURL = mcRouterHubCentralURL
	app.McRouterConfigFile = mcRouterConfigFile

	configManager, err := models.NewMcRouterConfigManager(app.McRouterConfigFile)
	if err != nil {
		return nil, err
	}
	app.McRouterConfigManager = configManager

	app.NrInsightsURL = nrInsightsURL
	app.NrInsightsInsertKey = nrInsightsInsertKey

	if app.ReportInterval == "" {
		app.ReportInterval = "3s"
	}

	app.Storage = storage.New()

	return app, nil
}

type Application struct {
	McRouterAddr          string
	McRouterConfigFile    string
	McRouterHubCentralURL string
	McRouterStatsManager  *models.McRouterStatsManager
	McRouterConfigManager *models.McRouterConfigManager
	NrInsightsURL         string
	NrInsightsInsertKey   string
	ReportInterval        string
	ReadOnly              bool
	Storage               *storage.Storage
}

func (app *Application) IsCentral() bool {
	return app.McRouterHubCentralURL == ""
}

func (app *Application) CollectData() error {
	go func() {
		for {
			stats, err := app.McRouterStatsManager.Stats()
			if err == nil && stats != nil {
				app.Storage.Set("stats", stats)
			}

			statsFromFile, err := app.McRouterStatsManager.StatsFromFile()
			if err == nil && stats != nil {
				app.Storage.Set("statsFromFile", statsFromFile)
			}

			config, err := app.McRouterConfigManager.Config()
			if err == nil && config != nil {
				app.Storage.Set("config", config)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	return nil
}

func (app *Application) ReportToCentral() error {
	if app.IsCentral() {
		return nil
	}

	duration, err := time.ParseDuration(app.ReportInterval)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	client := &http.Client{}

	url := app.McRouterHubCentralURL
	if !strings.HasSuffix(url, "/central") {
		url = url + "/central"
	}

	go func() {
		for {
			configInterface := app.Storage.Get("config")
			if configInterface == nil {
				time.Sleep(duration)
				continue
			}
			config := configInterface.(map[string]interface{})

			payload := &handlers.Payload{Hostname: hostname, Config: config}

			payloadJson, err := json.Marshal(payload)

			req, err := http.NewRequest("POST", url, bytes.NewReader(payloadJson))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to create HTTP request struct")

				time.Sleep(duration)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to send HTTP request")

				time.Sleep(duration)
				continue
			}

			defer resp.Body.Close()

			time.Sleep(duration)
		}
	}()

	return err
}

func (app *Application) ReportToNewrelicInsights() error {
	if app.NrInsightsInsertKey == "" {
		return nil
	}

	duration, err := time.ParseDuration(app.ReportInterval)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	client := &http.Client{}

	go func() {
		for {
			statsInterface := app.Storage.Get("stats")
			if statsInterface == nil {
				time.Sleep(duration)
				continue
			}
			stats := statsInterface.(*models.Stats)

			payload := structs.Map(stats)
			payload["hostname"] = hostname
			payload["eventType"] = "McRouter"

			// Fetch the other stats data from file.
			statsFromFileInterface := app.Storage.Get("statsFromFile")

			if statsFromFileInterface != nil {
				statsFromFile := statsFromFileInterface.(map[string]interface{})

				for key, value := range statsFromFile {
					trimmedKey := strings.Replace(key, "libmcrouter.mcrouter.5000.", "", -1)
					payload[trimmedKey] = value
				}
			}

			payloadJson, err := json.Marshal(payload)

			req, err := http.NewRequest("POST", app.NrInsightsURL, bytes.NewReader(payloadJson))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to create HTTP request struct")

				time.Sleep(duration)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Insert-Key", app.NrInsightsInsertKey)

			resp, err := client.Do(req)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to send HTTP request")

				time.Sleep(duration)
				continue
			}

			defer resp.Body.Close()

			time.Sleep(duration)
		}
	}()

	return err
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetMcRouterConfigFile(app.McRouterConfigFile))
	middle.Use(middlewares.SetStorage(app.Storage))
	middle.Use(middlewares.SetReadOnly(app.ReadOnly))
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetRoot).Methods("GET")
	router.HandleFunc("/config", handlers.GetConfig).Methods("GET")

	router.HandleFunc("/config/pools", handlers.GetConfigPools).Methods("GET")

	if app.IsCentral() {
		router.HandleFunc("/central", handlers.GetCentral).Methods("GET")
	}

	return app.addWriteHandlers(router)
}

func (app *Application) addWriteHandlers(router *mux.Router) *mux.Router {
	if !app.ReadOnly {
		router.HandleFunc("/config", handlers.PostConfig).Methods("POST")

		if app.IsCentral() {
			router.HandleFunc("/central", handlers.PostCentral).Methods("POST")
		}
	}

	return router
}
