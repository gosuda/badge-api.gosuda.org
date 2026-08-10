package main

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	if got := response.Header().Get("X-Robots-Tag"); got != "noindex, noarchive" {
		t.Fatalf("X-Robots-Tag = %q", got)
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

func TestHealthzIsNotIndexable(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("X-Robots-Tag"); got != "noindex, noarchive" {
		t.Fatalf("X-Robots-Tag = %q", got)
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

func TestBadgeSizeScalesEveryStyle(t *testing.T) {
	for _, style := range availableStyles() {
		t.Run(style, func(t *testing.T) {
			base := renderValidatedBadge(t, badgeOptions{Message: "ready", Style: style, Size: 100})
			scaled := renderValidatedBadge(t, badgeOptions{Message: "ready", Style: style, Size: 175})
			baseRoot := parseSVGRoot(t, base)
			scaledRoot := parseSVGRoot(t, scaled)

			baseWidth, err := strconv.ParseFloat(baseRoot.Width, 64)
			if err != nil {
				t.Fatal(err)
			}
			scaledWidth, err := strconv.ParseFloat(scaledRoot.Width, 64)
			if err != nil {
				t.Fatal(err)
			}
			if scaledWidth != baseWidth*1.75 {
				t.Fatalf("scaled width = %v, want %v", scaledWidth, baseWidth*1.75)
			}

			baseHeight, err := strconv.ParseFloat(baseRoot.Height, 64)
			if err != nil {
				t.Fatal(err)
			}
			scaledHeight, err := strconv.ParseFloat(scaledRoot.Height, 64)
			if err != nil {
				t.Fatal(err)
			}
			if scaledHeight != baseHeight*1.75 {
				t.Fatalf("scaled height = %v, want %v", scaledHeight, baseHeight*1.75)
			}
			if scaledRoot.ViewBox != baseRoot.ViewBox {
				t.Fatalf("scaled viewBox = %q, want %q", scaledRoot.ViewBox, baseRoot.ViewBox)
			}
		})
	}
}

func TestBadgeHandlerRejectsInvalidSize(t *testing.T) {
	for _, size := range []string{"", "0", "49", "301", "large"} {
		t.Run(size, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/badge.svg?message=ready&size="+size, nil)
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestBadgeHandlerAppliesExactSize(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/badge.svg?message=ready&style=flat&size=150", nil)
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	root := parseSVGRoot(t, response.Body.Bytes())
	if root.Height != "30" {
		t.Fatalf("height = %q, want 30", root.Height)
	}
}

func TestOldSchoolStyleUsesFixedPixelGeometry(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/badge.svg?label=pixel&message=button&style=old-school&labelColor=ff5a18&color=a8a979&size=100", nil)
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	root := parseSVGRoot(t, response.Body.Bytes())
	if root.Width != "80" || root.Height != "15" || root.ViewBox != "0 0 80 15" {
		t.Fatalf("root = width %q, height %q, viewBox %q; want 80, 15, 0 0 80 15", root.Width, root.Height, root.ViewBox)
	}
	body := response.Body.String()
	for _, expected := range []string{
		`shape-rendering="crispEdges"`,
		`<rect width="80" height="15" fill="#a5a5a5"/>`,
		`<rect x="1" y="1" width="78" height="13" fill="#fff"/>`,
		`y="2" width="30" height="11" fill="#ff5a18"`,
		`fill="#a8a979"`,
		`<text `,
		`text-rendering="geometricPrecision"`,
		`lengthAdjust="spacingAndGlyphs"`,
		`>PIXEL</text>`,
		`aria-label="PIXEL: BUTTON"`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("SVG does not contain %q", expected)
		}
	}
	if !strings.Contains(body, ">BUTTON</text>") {
		t.Fatal("old-school SVG does not render the full message as adaptive vector text")
	}

	scaled := renderValidatedBadge(t, badgeOptions{Label: "pixel", Message: "button", Style: "old-school", Size: 200})
	scaledRoot := parseSVGRoot(t, scaled)
	if scaledRoot.Width != "160" || scaledRoot.Height != "30" || scaledRoot.ViewBox != root.ViewBox {
		t.Fatalf("scaled root = width %q, height %q, viewBox %q; want 160, 30, %q", scaledRoot.Width, scaledRoot.Height, scaledRoot.ViewBox, root.ViewBox)
	}
}

func TestClassicPixelStylesUseFixed88x31Geometry(t *testing.T) {
	tests := []struct {
		style   string
		label   string
		message string
		marker  string
	}{
		{style: "click-here", label: "click", message: "here", marker: `fill="#ea3323"`},
		{style: "best-viewed", label: "viewed with", message: "chrome", marker: `fill="#64a8d2"`},
	}

	for _, test := range tests {
		t.Run(test.style, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/badge.svg?label="+strings.ReplaceAll(test.label, " ", "%20")+"&message="+test.message+"&style="+test.style+"&size=100", nil)
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			root := parseSVGRoot(t, response.Body.Bytes())
			if root.Width != "88" || root.Height != "31" || root.ViewBox != "0 0 88 31" {
				t.Fatalf("root = width %q, height %q, viewBox %q; want 88, 31, 0 0 88 31", root.Width, root.Height, root.ViewBox)
			}
			if body := response.Body.String(); !strings.Contains(body, test.marker) || strings.Contains(body, "<text") {
				t.Fatalf("SVG does not preserve expected pixel artwork marker %q", test.marker)
			}

			scaled := renderValidatedBadge(t, badgeOptions{Label: test.label, Message: test.message, Style: test.style, Size: 200})
			scaledRoot := parseSVGRoot(t, scaled)
			if scaledRoot.Width != "176" || scaledRoot.Height != "62" || scaledRoot.ViewBox != root.ViewBox {
				t.Fatalf("scaled root = width %q, height %q, viewBox %q; want 176, 62, %q", scaledRoot.Width, scaledRoot.Height, scaledRoot.ViewBox, root.ViewBox)
			}
		})
	}
}

type svgRoot struct {
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
}

func renderValidatedBadge(t *testing.T, options badgeOptions) []byte {
	t.Helper()
	validated, err := validateBadgeOptions(options)
	if err != nil {
		t.Fatal(err)
	}
	return renderBadge(validated)
}

func parseSVGRoot(t *testing.T, svg []byte) svgRoot {
	t.Helper()
	var root svgRoot
	if err := xml.Unmarshal(svg, &root); err != nil {
		t.Fatal(err)
	}
	return root
}
