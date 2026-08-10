package main

import (
	"net/http"
	"net/http/cgi"
	"os"
	"strings"
)

const gatewayInterfaceEnvironment = "GATEWAY_INTERFACE"

func cgiModeEnabled() bool {
	return os.Getenv(gatewayInterfaceEnvironment) != ""
}

func serveCGI(handler http.Handler) error {
	return cgi.Serve(stripCGIScriptName(handler, os.Getenv("SCRIPT_NAME")))
}

func stripCGIScriptName(handler http.Handler, scriptName string) http.Handler {
	prefix := strings.TrimSuffix(scriptName, "/")
	if prefix == "" {
		return handler
	}

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if path != prefix && !strings.HasPrefix(path, prefix+"/") {
			handler.ServeHTTP(response, request)
			return
		}

		requestCopy := request.Clone(request.Context())
		urlCopy := *request.URL
		urlCopy.Path = strings.TrimPrefix(path, prefix)
		if urlCopy.Path == "" {
			urlCopy.Path = "/"
		}
		if urlCopy.RawPath != "" {
			urlCopy.RawPath = strings.TrimPrefix(urlCopy.RawPath, prefix)
			if urlCopy.RawPath == "" {
				urlCopy.RawPath = "/"
			}
		}
		requestCopy.URL = &urlCopy
		handler.ServeHTTP(response, requestCopy)
	})
}
