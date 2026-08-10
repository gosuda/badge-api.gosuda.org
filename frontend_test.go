package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestFrontendFileHandlerMissingAssetsReturnsNotFound(t *testing.T) {
	files := fstest.MapFS{
		"dist/.gitkeep": new(fstest.MapFile),
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	frontendFileHandler(files).ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", got)
	}
}

func TestFrontendFileHandlerServesBuiltIndex(t *testing.T) {
	files := fstest.MapFS{
		"dist/frontend/index.html": {Data: []byte("<!doctype html><title>Tiny Badge</title>")},
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	frontendFileHandler(files).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", got)
	}
}

func TestFrontendFileHandlerDoesNotCacheDirectoryIndex(t *testing.T) {
	files := fstest.MapFS{
		"dist/frontend/pricing/index.html": {Data: []byte("<!doctype html><title>Pricing</title>")},
	}
	request := httptest.NewRequest(http.MethodGet, "/pricing/", nil)
	response := httptest.NewRecorder()

	frontendFileHandler(files).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", got)
	}
}
