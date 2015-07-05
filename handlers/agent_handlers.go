package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/didip/mcrouter-hub/libhttp"
	"github.com/didip/mcrouter-hub/models"
	"github.com/fatih/structs"
	"github.com/gorilla/context"
)

func AgentGetRoot(w http.ResponseWriter, r *http.Request) {
	readOnly := context.Get(r, "readOnly").(bool)
	if readOnly {
		w.Write([]byte(`{
    paths: {
        GET: [
            "/config",
            "/config/pools",
            "/stats"
        ]
    }
}`))

	} else {
		w.Write([]byte(`{
    paths: {
        GET: [
            "/config",
            "/config/pools",
            "/stats"
        ],
        POST: [
            "/config"
        ]
    }
}`))
	}

}

func AgentGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mcRouterConfigFile := context.Get(r, "mcRouterConfigFile").(string)
	if mcRouterConfigFile == "" {
		err := errors.New("McRouter config file is missing")
		libhttp.HandleErrorJson(w, err)
		return
	}

	configManager, err := models.NewMcRouterConfigManager(mcRouterConfigFile)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	mcRouterConfigJson, err := configManager.ConfigJson()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(mcRouterConfigJson)
}

func AgentPostConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mcRouterConfigFile := context.Get(r, "mcRouterConfigFile").(string)
	if mcRouterConfigFile == "" {
		err := errors.New("McRouter config file is missing")
		libhttp.HandleErrorJson(w, err)
		return
	}

	configManager, err := models.NewMcRouterConfigManager(mcRouterConfigFile)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	mcRouterConfigJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = configManager.UpdateConfigJson(mcRouterConfigJson)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	libhttp.HandleSuccessJson(w, "New config is saved successfully")
}

func AgentGetConfigPools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mcRouterConfigFile := context.Get(r, "mcRouterConfigFile").(string)
	if mcRouterConfigFile == "" {
		err := errors.New("McRouter config file is missing")
		libhttp.HandleErrorJson(w, err)
		return
	}

	configManager, err := models.NewMcRouterConfigManager(mcRouterConfigFile)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	poolsJson, err := configManager.PoolsJson()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(poolsJson)
}

func AgentPostConfigPools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mcRouterConfigFile := context.Get(r, "mcRouterConfigFile").(string)
	if mcRouterConfigFile == "" {
		err := errors.New("McRouter config file is missing")
		libhttp.HandleErrorJson(w, err)
		return
	}

	configManager, err := models.NewMcRouterConfigManager(mcRouterConfigFile)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	poolsJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = configManager.UpdatePoolsJson(poolsJson)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	libhttp.HandleSuccessJson(w, "New config is saved successfully")
}

func AgentGetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statsInterface := context.Get(r, "stats")
	if statsInterface == nil {
		w.Write([]byte(`{}`))
		return
	}

	stats := statsInterface.(*models.Stats)
	payload := structs.Map(stats)

	// Fetch the other stats data from file.
	statsFromFileInterface := context.Get(r, "statsFromFile")

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

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(payloadJson)
}
