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
            "/configs"
        ]
    }
}`))

	} else {
		w.Write([]byte(`{
    paths: {
        GET: [
            "/configs"
        ],
        POST: [
            "/configs"
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

	store.Set(payload.Hostname, payload.Config)

	libhttp.HandleSuccessJson(w, "Config on host: "+payload.Hostname+" is saved successfully")
}

func CentralGetConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	allConfig, err := store.ToJson()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(allConfig)
}
