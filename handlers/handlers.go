package handlers

import (
	"errors"
	"github.com/didip/mcrouter-hub/libhttp"
	"github.com/didip/mcrouter-hub/models"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{
    paths: {
        GET: [
            "/config",
            "/config/pools"
        ],
        POST: [
            "/config"
        ]
    }
}`))
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
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

func PostConfig(w http.ResponseWriter, r *http.Request) {
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

func GetConfigPools(w http.ResponseWriter, r *http.Request) {
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

func PostConfigPools(w http.ResponseWriter, r *http.Request) {
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
