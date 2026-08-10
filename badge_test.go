package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBadgeHandlerCachesImmutableSVG(t *testing.T) {
	handler := newHandler()
	request := httptest.NewRequest(http.MethodGet, "/badge.svg?label=build&message=passing&style=flatbar&color=7c5cff", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Cache-Control"); got != "public, max-age=315360000, immutable" {
		t.Fatalf("Cache-Control = %q", got)
	}
	if got := response.Header().Get("Content-Type"); got != "image/svg+xml; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	etag := response.Header().Get("ETag")
	if etag == "" {
		t.Fatal("ETag is empty")
	}
	body := response.Body.String()
	for _, expected := range []string{`height="28"`, `BUILD`, `PASSING`, `fill="#7c5cff"`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("SVG does not contain %q", expected)
		}
	}

	conditional := httptest.NewRequest(http.MethodGet, request.URL.String(), nil)
	conditional.Header.Set("If-None-Match", etag)
	conditionalResponse := httptest.NewRecorder()
	handler.ServeHTTP(conditionalResponse, conditional)
	if conditionalResponse.Code != http.StatusNotModified {
		t.Fatalf("conditional status = %d, want %d", conditionalResponse.Code, http.StatusNotModified)
	}
	if conditionalResponse.Body.Len() != 0 {
		t.Fatalf("conditional body length = %d, want 0", conditionalResponse.Body.Len())
	}
}

func TestBadgeHandlerRejectsInvalidColor(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/badge.svg?message=test&color=not-a-color", nil)
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
}

func TestQueryMessagePreservesSVGSuffix(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		message string
	}{
		{name: "query", target: "/badge.svg?message=version.svg", message: "version.svg"},
		{name: "path", target: "/badge/release/version.svg", message: "version"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.target, nil)
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			if !strings.Contains(response.Body.String(), `aria-label="`+test.message+`"`) &&
				!strings.Contains(response.Body.String(), `aria-label="release: `+test.message+`"`) {
				t.Fatalf("SVG does not preserve message %q", test.message)
			}
		})
	}
}

func TestEveryStyleRendersEscapedSVG(t *testing.T) {
	for _, style := range availableStyles() {
		t.Run(style, func(t *testing.T) {
			options, err := validateBadgeOptions(badgeOptions{
				Label:            "<build>",
				Message:          "safe & stable",
				Style:            style,
				LabelColor:       "292724",
				MessageColor:     "d6ef53",
				LabelTextColor:   "ffffff",
				MessageTextColor: "292724",
			})
			if err != nil {
				t.Fatal(err)
			}
			svg := string(renderBadge(options))
			if !strings.HasPrefix(svg, "<svg ") {
				t.Fatalf("output does not start with SVG: %q", svg[:min(len(svg), 40)])
			}
			if strings.Contains(svg, "<build>") || strings.Contains(svg, "safe & stable") {
				t.Fatal("SVG contains unescaped text")
			}
			if !strings.Contains(svg, "&lt;") || !strings.Contains(svg, "&amp;") {
				t.Fatal("SVG does not contain escaped text")
			}
		})
	}
}
