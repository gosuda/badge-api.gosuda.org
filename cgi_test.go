package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCGIModeEnabled(t *testing.T) {
	setOptionalEnvironment(t, gatewayInterfaceEnvironment, nil)
	if cgiModeEnabled() {
		t.Fatal("CGI mode enabled without GATEWAY_INTERFACE")
	}

	setOptionalEnvironment(t, gatewayInterfaceEnvironment, new("CGI/1.1"))
	if !cgiModeEnabled() {
		t.Fatal("CGI mode disabled with GATEWAY_INTERFACE")
	}
}

func TestStripCGIScriptName(t *testing.T) {
	tests := []struct {
		name       string
		requestURI string
		wantPath   string
	}{
		{name: "root", requestURI: "/cgi-bin/badge-api.gosuda.org.cgi", wantPath: "/"},
		{name: "badge", requestURI: "/cgi-bin/badge-api.gosuda.org.cgi/badge.svg?message=ready", wantPath: "/badge.svg"},
		{name: "pricing", requestURI: "/cgi-bin/badge-api.gosuda.org.cgi/pricing/", wantPath: "/pricing/"},
		{name: "already stripped", requestURI: "/reference/", wantPath: "/reference/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotPath string
			handler := stripCGIScriptName(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				gotPath = request.URL.Path
				response.WriteHeader(http.StatusNoContent)
			}), "/cgi-bin/badge-api.gosuda.org.cgi")
			request := httptest.NewRequest(http.MethodGet, test.requestURI, nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
			}
			if gotPath != test.wantPath {
				t.Fatalf("path = %q, want %q", gotPath, test.wantPath)
			}
		})
	}
}
