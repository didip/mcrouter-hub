package handlers

import (
	"errors"
	"github.com/didip/mcrouter-hub/libhttp"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"os"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{
    paths: {
        GET: [
            "/config"
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

	mcRouterConfigJson, err := ioutil.ReadFile(mcRouterConfigFile)
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

	fileInfo, err := os.Stat(mcRouterConfigFile)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	mcRouterConfigJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	err = ioutil.WriteFile(mcRouterConfigFile, mcRouterConfigJson, fileInfo.Mode())
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	libhttp.HandleSuccessJson(w, "New config is saved successfully")
}
