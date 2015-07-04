// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"github.com/didip/mcrouter-hub/storage"
	"github.com/gorilla/context"
)

func SetReadOnly(readOnly bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "readOnly", readOnly)

			next.ServeHTTP(res, req)
		})
	}
}

func SetMcRouterConfigFile(McRouterConfigFile string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "mcRouterConfigFile", McRouterConfigFile)

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
