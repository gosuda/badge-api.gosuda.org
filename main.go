package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

const badgeCacheSeconds = 315360000

//go:embed all:dist
var distributionFiles embed.FS

func main() {
	if cgiModeEnabled() {
		if err := serveCGI(newHandler()); err != nil {
			log.Fatal(err)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	exitOnStdinEOF, err := stdinEOFShutdownEnabled()
	if err != nil {
		log.Fatal(err)
	}
	if exitOnStdinEOF {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		go shutdownOnStdinEOF(ctx, os.Stdin, cancel)
	}

	target, err := configuredListenerTarget()
	if err != nil {
		log.Fatal(err)
	}
	if err := removeUnixSocket(target); err != nil {
		log.Fatal(err)
	}
	defer removeUnixSocket(target)

	listener, err := net.Listen(target.Network, target.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	if err := secureUnixSocket(target); err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              target.Address,
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownResult := make(chan error, 1)
	go func() {
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownResult <- server.Shutdown(shutdownContext)
	}()

	log.Printf("badge API listening on %s %s", target.Network, listener.Addr())
	serveErr := server.Serve(listener)
	if ctx.Err() != nil {
		if shutdownErr := <-shutdownResult; shutdownErr != nil {
			log.Printf("graceful shutdown failed: %v", shutdownErr)
		}
	}
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		log.Fatal(serveErr)
	}
}

func newHandler() http.Handler {
	router := httprouter.New()
	router.GET("/badge.svg", serveBadge)
	router.HEAD("/badge.svg", serveBadge)
	router.GET("/badge/:label/:message", serveBadge)
	router.HEAD("/badge/:label/:message", serveBadge)
	router.GET("/healthz", func(response http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		response.Header().Set("Cache-Control", "no-store")
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write([]byte("ok\n"))
	})

	router.NotFound = frontendFileHandler(distributionFiles)
	return requestHeaders(router)
}

func frontendFileHandler(files fs.FS) http.Handler {
	if _, err := fs.Stat(files, "dist/frontend"); err != nil {
		return frontendHeaders(http.NotFoundHandler())
	}
	frontend, err := fs.Sub(files, "dist/frontend")
	if err != nil {
		return frontendHeaders(http.NotFoundHandler())
	}
	return frontendHeaders(http.FileServer(http.FS(frontend)))
}

func serveBadge(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	options, err := badgeOptionsFromRequest(request, params)
	if err != nil {
		response.Header().Set("Cache-Control", "no-store")
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	svg := renderBadge(options)
	hash := sha256.Sum256(svg)
	etag := `"` + hex.EncodeToString(hash[:16]) + `"`
	setBadgeHeaders(response.Header(), etag)
	if request.Header.Get("If-None-Match") == etag {
		response.WriteHeader(http.StatusNotModified)
		return
	}
	response.Header().Set("Content-Length", fmt.Sprintf("%d", len(svg)))
	response.WriteHeader(http.StatusOK)
	if request.Method != http.MethodHead {
		_, _ = response.Write(svg)
	}
}

func badgeOptionsFromRequest(request *http.Request, params httprouter.Params) (badgeOptions, error) {
	query := request.URL.Query()
	label := params.ByName("label")
	pathMessage := params.ByName("message")
	message := strings.TrimSuffix(pathMessage, ".svg")
	if query.Has("label") {
		label = query.Get("label")
	}
	if query.Has("message") {
		message = query.Get("message")
	}
	if label == "_" {
		label = ""
	}

	return validateBadgeOptions(badgeOptions{
		Label:            label,
		Message:          message,
		Style:            query.Get("style"),
		LabelColor:       query.Get("labelColor"),
		MessageColor:     query.Get("color"),
		LabelTextColor:   query.Get("labelTextColor"),
		MessageTextColor: query.Get("textColor"),
	})
}

func setBadgeHeaders(header http.Header, etag string) {
	cacheValue := fmt.Sprintf("public, max-age=%d, immutable", badgeCacheSeconds)
	header.Set("Cache-Control", cacheValue)
	header.Set("CDN-Cache-Control", cacheValue)
	header.Set("Surrogate-Control", cacheValue)
	header.Set("Content-Type", "image/svg+xml; charset=utf-8")
	header.Set("Content-Disposition", "inline")
	header.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("ETag", etag)
}

func frontendHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		switch {
		case path == "/" || strings.HasSuffix(path, ".html"):
			response.Header().Set("Cache-Control", "no-cache")
		case strings.HasPrefix(path, "/_app/immutable/"):
			response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		default:
			response.Header().Set("Cache-Control", "public, max-age=3600")
		}
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		response.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		response.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; connect-src 'self'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'")
		next.ServeHTTP(response, request)
	})
}

func requestHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(response, request)
	})
}
