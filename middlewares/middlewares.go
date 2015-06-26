// Package middlewares provides common middleware handlers.
package middlewares

import (
	"github.com/didip/mcrouter-hub/storage"
	"github.com/gorilla/context"
	"net/http"
)

func SetMcRounterConfigFile(McRounterConfigFile string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "mcRouterConfigFile", McRounterConfigFile)

			next.ServeHTTP(res, req)
		})
	}
}

func SetStorage(store *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "store", store)

			next.ServeHTTP(res, req)
		})
	}
}
