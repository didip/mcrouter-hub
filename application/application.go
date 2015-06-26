package application

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/didip/mcrouter-hub/handlers"
	"github.com/didip/mcrouter-hub/middlewares"
	"github.com/didip/mcrouter-hub/models"
	"github.com/didip/mcrouter-hub/storage"
	"github.com/go-fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
	"time"
)

func New(mcRouterConfigFile, mcRouterHubCentralURL string) (*Application, error) {
	app := &Application{}
	app.McRounterConfigFile = mcRouterConfigFile
	app.McRouterHubCentralURL = mcRouterHubCentralURL

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if app.ReportInterval == "" {
		app.ReportInterval = "1m"
	}

	app.Storage = storage.New()
	app.fileWatcher = watcher

	return app, nil
}

type Application struct {
	McRounterConfigFile   string
	McRouterHubCentralURL string
	ReportInterval        string
	Storage               *storage.Storage
	fileWatcher           *fsnotify.Watcher
}

func (app *Application) IsCentral() bool {
	return app.McRouterHubCentralURL == ""
}

func (app *Application) WatchMcRounterConfigFile() error {
	defer app.fileWatcher.Close()

	go func() {
		for {
			select {
			case event := <-app.fileWatcher.Events:
				if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Remove == fsnotify.Remove) || (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) {
					logrus.WithFields(logrus.Fields{
						"Event": event.String(),
						"File":  event.Name,
					}).Info("Config file changed")
				}
			case err := <-app.fileWatcher.Errors:
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"Error": err.Error(),
					}).Error("Error while watching config file")
				}
			}
		}
	}()

	err := app.fileWatcher.Add(app.McRounterConfigFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
			"File":  app.McRounterConfigFile,
		}).Error("Error while adding config file to watch")
	}

	return err
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

	configManager, err := models.NewMcRouterConfigManager(app.McRounterConfigFile)
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
			config, err := configManager.Config()
			if err != nil {
				time.Sleep(duration)
				continue
			}

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

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetMcRounterConfigFile(app.McRounterConfigFile))
	middle.Use(middlewares.SetStorage(app.Storage))
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetRoot).Methods("GET")
	router.HandleFunc("/config", handlers.GetConfig).Methods("GET")
	router.HandleFunc("/config", handlers.PostConfig).Methods("POST")

	router.HandleFunc("/config/pools", handlers.GetConfigPools).Methods("GET")

	if app.IsCentral() {
		router.HandleFunc("/central", handlers.PostCentral).Methods("POST")
		router.HandleFunc("/central", handlers.GetCentral).Methods("GET")
	}

	return router
}
