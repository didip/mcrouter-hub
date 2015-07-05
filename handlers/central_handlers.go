package handlers

import (
	"encoding/json"
	"github.com/didip/mcrouter-hub/libhttp"
	"github.com/didip/mcrouter-hub/payloads"
	"github.com/didip/mcrouter-hub/storage"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
)

func CentralGetRoot(w http.ResponseWriter, r *http.Request) {
	readOnly := context.Get(r, "readOnly").(bool)
	if readOnly {
		w.Write([]byte(`{
    paths: {
        GET: [
            "/configs",
            "/stats"
        ]
    }
}`))

	} else {
		w.Write([]byte(`{
    paths: {
        GET: [
            "/configs",
            "/stats"
        ],
        POST: [
            "/configs",
            "/stats"
        ]
    }
}`))
	}

}

func CentralPostConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	payloadJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	var payload payloads.ReportConfigToCentralPayload
	err = json.Unmarshal(payloadJson, &payload)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	store.Set("config:"+payload.Hostname, payload.Config)

	libhttp.HandleSuccessJson(w, "Config on host: "+payload.Hostname+" is saved successfully")
}

func CentralGetConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	configs, err := store.ToJson("config:")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(configs)
}

func CentralPostStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	payloadJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	payload := make(map[string]interface{})
	err = json.Unmarshal(payloadJson, &payload)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	hostname := ""
	hostnameInterface := payload["hostname"]
	if hostnameInterface != nil {
		hostname = hostnameInterface.(string)
	}

	store.Set("stats:"+hostname, payload)

	libhttp.HandleSuccessJson(w, "Stats on host: "+hostname+" is saved successfully")
}

func CentralGetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	stats, err := store.ToJson("stats:")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(stats)
}
