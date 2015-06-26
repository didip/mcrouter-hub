// Package libhttp provides http related library functions.
package libhttp

import (
	"encoding/json"
	"net/http"
)

// HandleErrorJson wraps error in JSON structure.
func HandleErrorJson(w http.ResponseWriter, err error) {
	var errMap map[string]string

	if err == nil {
		errMap = map[string]string{"Error": "Error struct is nil."}
	} else {
		errMap = map[string]string{"Error": err.Error()}
	}

	errJson, _ := json.Marshal(errMap)
	http.Error(w, string(errJson), http.StatusInternalServerError)
}

func HandleSuccessJson(w http.ResponseWriter, message string) {
	jsonMap := map[string]string{"Success": message}

	successJSON, _ := json.Marshal(jsonMap)
	w.Write(successJSON)
}
