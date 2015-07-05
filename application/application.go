package application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/didip/mcrouter-hub/handlers"
	"github.com/didip/mcrouter-hub/middlewares"
	"github.com/didip/mcrouter-hub/models"
	"github.com/didip/mcrouter-hub/payloads"
	"github.com/didip/mcrouter-hub/storage"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
)

func New() (*Application, error) {
	app := &Application{}
	app.Settings = make(map[string]string)
	app.Storage = storage.New()

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		app.Settings[pair[0]] = pair[1]
	}

	if app.Settings["MCRHUB_MODE"] == "" {
		app.Settings["MCRHUB_MODE"] = "agent"
	}
	if app.Settings["MCRHUB_MODE"] == "agent" && app.Settings["MCROUTER_ADDR"] == "" {
		return nil, errors.New("MCROUTER_ADDR is required")
	}
	if app.Settings["MCRHUB_MODE"] == "agent" && app.Settings["MCROUTER_CONFIG_FILE"] == "" {
		return nil, errors.New("MCROUTER_CONFIG_FILE is required")
	}
	if app.Settings["MCRHUB_MODE"] == "agent" && app.Settings["MCRHUB_READ_ONLY"] == "" {
		app.Settings["MCRHUB_READ_ONLY"] = "true"
	}
	if app.Settings["MCRHUB_MODE"] == "central" && app.Settings["MCRHUB_READ_ONLY"] == "" {
		app.Settings["MCRHUB_READ_ONLY"] = "false"
	}
	if app.Settings["MCRHUB_REPORT_INTERVAL"] == "" {
		app.Settings["MCRHUB_REPORT_INTERVAL"] = "3s"
	}
	if app.Settings["MCRHUB_ADDR"] == "" {
		if app.IsAgentMode() {
			app.Settings["MCRHUB_ADDR"] = ":5001"
		} else {
			app.Settings["MCRHUB_ADDR"] = ":5002"
		}
	}

	if app.IsAgentMode() {
		app.McRouterStatsManager = models.NewMcRouterStatsManager(app.Settings["MCROUTER_ADDR"])

		configManager, err := models.NewMcRouterConfigManager(app.Settings["MCROUTER_CONFIG_FILE"])
		if err != nil {
			return nil, err
		}
		app.McRouterConfigManager = configManager
	}

	return app, nil
}

type Application struct {
	Settings              map[string]string
	McRouterStatsManager  *models.McRouterStatsManager
	McRouterConfigManager *models.McRouterConfigManager
	Storage               *storage.Storage
}

func (app *Application) SettingKeys() []string {
	return []string{
		"MCROUTER_ADDR",
		"MCROUTER_CONFIG_FILE",
		"MCRHUB_MODE",
		"MCRHUB_READ_ONLY",
		"MCRHUB_CENTRAL_URLS",
		"MCRHUB_LOG_LEVEL",
		"MCRHUB_REPORT_INTERVAL",
		"MCRHUB_TOKENS_DIR",
		"MCRHUB_ADDR",
		"MCRHUB_CERT_FILE",
		"MCRHUB_KEY_FILE",
		"NR_INSIGHTS_URL",
		"NR_INSIGHTS_INSERT_KEY",
	}
}

func (app *Application) IsReadOnly() bool {
	if strings.ToLower(app.Settings["MCRHUB_READ_ONLY"]) == "false" {
		return false
	}
	return true
}

func (app *Application) IsAgentMode() bool {
	if strings.ToLower(app.Settings["MCRHUB_MODE"]) == "agent" {
		return true
	}
	return false
}

func (app *Application) IsCentralMode() bool {
	if strings.ToLower(app.Settings["MCRHUB_MODE"]) == "central" {
		return true
	}
	return false
}

func (app *Application) CentralURLs() []string {
	urls := make([]string, 0)

	if app.Settings["MCRHUB_CENTRAL_URLS"] == "" {
		return urls
	}

	urlParts := strings.Split(app.Settings["MCRHUB_CENTRAL_URLS"], ",")
	for _, urlPart := range urlParts {
		urlPartTrimmed := strings.TrimSpace(urlPart)
		urls = append(urls, urlPartTrimmed)
	}

	return urls
}

func (app *Application) Tokens() []string {
	tokens := make([]string, 0)

	if app.Settings["MCRHUB_TOKENS_DIR"] == "" {
		return tokens
	}

	app.Settings["MCRHUB_TOKENS_DIR"] = os.ExpandEnv(app.Settings["MCRHUB_TOKENS_DIR"])

	fileInfos, err := ioutil.ReadDir(app.Settings["MCRHUB_TOKENS_DIR"])
	if err != nil {
		return tokens
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			fullpath := filepath.Join(app.Settings["MCRHUB_TOKENS_DIR"], fileInfo.Name())

			file, err := os.Open(fullpath)
			if err != nil {
				continue
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				token := scanner.Text()
				if token != "" {
					tokens = append(tokens, token)
				}
			}
		}
	}

	return tokens
}

func (app *Application) GetStats() map[string]interface{} {
	statsInterface := app.Storage.Get("stats")
	if statsInterface == nil {
		return nil
	}

	stats := statsInterface.(*models.Stats)
	payload := structs.Map(stats)

	// Fetch the other stats data from file.
	statsFromFileInterface := app.Storage.Get("statsFromFile")

	if statsFromFileInterface != nil {
		statsFromFile := statsFromFileInterface.(map[string]interface{})

		for key, value := range statsFromFile {
			trimmedKey := strings.Replace(key, "libmcrouter.mcrouter.5000.", "", -1)
			payload[trimmedKey] = value
		}
	}

	hostname, err := os.Hostname()
	if err == nil {
		payload["hostname"] = hostname
	}

	return payload
}

func (app *Application) CollectData() error {
	if !app.IsAgentMode() {
		return nil
	}

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

			statsPayload := app.GetStats()
			app.Storage.Set("statsPayload", statsPayload)

			config, err := app.McRouterConfigManager.Config()
			if err == nil && config != nil {
				app.Storage.Set("config", config)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	return nil
}

func (app *Application) ReportConfigToCentral() error {
	if !app.IsAgentMode() {
		return nil
	}
	if len(app.CentralURLs()) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(app.Settings["MCRHUB_REPORT_INTERVAL"])
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	client := &http.Client{}

	for _, url := range app.CentralURLs() {
		if !strings.HasSuffix(url, "/configs") {
			url = url + "/configs"
		}

		go func() {
			for {
				configInterface := app.Storage.Get("config")
				if configInterface == nil {
					time.Sleep(duration)
					continue
				}
				config := configInterface.(map[string]interface{})

				payload := &payloads.ReportConfigToCentralPayload{Hostname: hostname, Config: config}

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
	}

	return nil
}

func (app *Application) ReportStatsToCentral() error {
	if !app.IsAgentMode() {
		return nil
	}
	if len(app.CentralURLs()) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(app.Settings["MCRHUB_REPORT_INTERVAL"])
	if err != nil {
		return err
	}

	client := &http.Client{}

	for _, url := range app.CentralURLs() {
		if !strings.HasSuffix(url, "/stats") {
			url = url + "/stats"
		}

		go func() {
			for {
				payload := app.GetStats()
				if payload == nil {
					logrus.Error("Failed to get stats")
					time.Sleep(duration)
					continue
				}

				payloadJson, err := json.Marshal(payload)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"Error": err.Error(),
					}).Error("Failed to create marshal JSON")

					time.Sleep(duration)
					continue
				}

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
	}

	return nil
}

func (app *Application) ReportStatsToNewrelicInsights() error {
	if app.Settings["NR_INSIGHTS_URL"] == "" || app.Settings["NR_INSIGHTS_INSERT_KEY"] == "" {
		return nil
	}

	duration, err := time.ParseDuration(app.Settings["MCRHUB_REPORT_INTERVAL"])
	if err != nil {
		return err
	}

	client := &http.Client{}

	go func() {
		for {
			payload := app.GetStats()
			if payload == nil {
				logrus.Error("Failed to get stats")
				time.Sleep(duration)
				continue
			}

			payload["eventType"] = "McRouter"

			payloadJson, err := json.Marshal(payload)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to create marshal JSON")

				time.Sleep(duration)
				continue
			}

			req, err := http.NewRequest("POST", app.Settings["NR_INSIGHTS_URL"], bytes.NewReader(payloadJson))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error": err.Error(),
				}).Error("Failed to create HTTP request struct")

				time.Sleep(duration)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Insert-Key", app.Settings["NR_INSIGHTS_INSERT_KEY"])

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
	middle.Use(middlewares.SetMcRouterConfigFile(app.Settings["MCROUTER_CONFIG_FILE"]))
	middle.Use(middlewares.SetStorage(app.Storage))
	middle.Use(middlewares.SetReadOnly(app.IsReadOnly()))
	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *mux.Router {
	router := mux.NewRouter()
	router = app.addAgentHandlers(router)
	router = app.addCentralHandlers(router)

	return router
}

func (app *Application) addAgentHandlers(router *mux.Router) *mux.Router {
	if app.IsAgentMode() {
		router.HandleFunc("/", handlers.AgentGetRoot).Methods("GET")
		router.HandleFunc("/config", handlers.AgentGetConfig).Methods("GET")
		router.HandleFunc("/config/pools", handlers.AgentGetConfigPools).Methods("GET")
		router.HandleFunc("/stats", handlers.AgentGetStats).Methods("GET")

		if !app.IsReadOnly() {
			router.HandleFunc("/config", handlers.AgentPostConfig).Methods("POST")
		}
	}
	return router
}

func (app *Application) addCentralHandlers(router *mux.Router) *mux.Router {
	if app.IsCentralMode() {
		router.HandleFunc("/", handlers.CentralGetRoot).Methods("GET")
		router.HandleFunc("/configs", handlers.CentralGetConfigs).Methods("GET")
		router.HandleFunc("/configs/{hostname}", handlers.CentralGetConfigsHostname).Methods("GET")
		router.HandleFunc("/stats", handlers.CentralGetStats).Methods("GET")
		router.HandleFunc("/stats/{hostname}", handlers.CentralGetStatsHostname).Methods("GET")

		if !app.IsReadOnly() {
			router.HandleFunc("/configs", handlers.CentralPostConfigs).Methods("POST")
			router.HandleFunc("/stats", handlers.CentralPostStats).Methods("POST")
		}
	}
	return router
}
