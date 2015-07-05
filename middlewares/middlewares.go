// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"github.com/didip/mcrouter-hub/libhttp"
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

// MustLoginApi is a middleware that checks POST with valid tokens.
func MustLoginApi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		store := context.Get(req, "store").(*storage.Storage)

		tokens := store.All("tokens:")

		auth := req.Header.Get("Authorization")

		// Ignore authentication if we have 0 tokens.
		if auth == "" && len(tokens) == 0 {
			next.ServeHTTP(res, req)
			return
		}

		if auth == "" {
			libhttp.BasicAuthUnauthorized(res, nil)
			return
		}

		accessTokenString, _, ok := libhttp.ParseBasicAuth(auth)
		if !ok {
			libhttp.BasicAuthUnauthorized(res, nil)
			return
		}

		tokenFound := false

		if len(tokens) == 0 {
			tokenFound = true
		}

		for token, _ := range tokens {
			if token == accessTokenString {
				tokenFound = true
			}
		}

		if !tokenFound {
			libhttp.BasicAuthUnauthorized(res, nil)
			return
		}

		next.ServeHTTP(res, req)
	})
}
