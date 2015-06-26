package handlers

import (
	"encoding/json"
	"github.com/didip/mcrouter-hub/libhttp"
	"github.com/didip/mcrouter-hub/storage"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	Hostname string
	Config   map[string]interface{}
}

func PostCentral(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	payloadJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	var payload Payload
	err = json.Unmarshal(payloadJson, &payload)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	store.Set(payload.Hostname, payload.Config)

	libhttp.HandleSuccessJson(w, "Config on host: "+payload.Hostname+" is saved successfully")
}

func GetCentral(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	store := context.Get(r, "store").(*storage.Storage)

	allConfig, err := store.ToJson()
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.Write(allConfig)
}
