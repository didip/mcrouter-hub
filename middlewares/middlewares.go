// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"github.com/gorilla/context"
)

func SetMcRounterConfigFile(McRounterConfigFile string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "mcRouterConfigFile", McRounterConfigFile)

			next.ServeHTTP(res, req)
		})
	}
}
