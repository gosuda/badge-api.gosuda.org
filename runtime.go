package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"badge-api.gosuda.org/internal/envaddr"
)

const (
	httpAddressEnvironment    = "BADGE_API_HTTP_ADDR"
	exitOnStdinEOFEnvironment = "BADGE_API_EXIT_ON_STDIN_EOF"
)

type listenerTarget struct {
	Network string
	Address string
}

func configuredListenerTarget() (listenerTarget, error) {
	value, configured := os.LookupEnv(httpAddressEnvironment)
	if !configured || value == "" {
		return listenerTarget{Network: "tcp", Address: envaddr.Get(":8080")}, nil
	}

	network, address, hasNetwork := strings.Cut(value, ":")
	if !hasNetwork || (network != "tcp" && network != "unix") {
		return listenerTarget{Network: "tcp", Address: value}, nil
	}
	if address == "" {
		return listenerTarget{}, fmt.Errorf("%s must include an address after %s:", httpAddressEnvironment, network)
	}
	return listenerTarget{Network: network, Address: address}, nil
}

func stdinEOFShutdownEnabled() (bool, error) {
	value, configured := os.LookupEnv(exitOnStdinEOFEnvironment)
	if !configured {
		return false, nil
	}
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s: %w", exitOnStdinEOFEnvironment, err)
	}
	return enabled, nil
}

func shutdownOnStdinEOF(ctx context.Context, reader io.Reader, cancel context.CancelFunc) {
	_, err := io.Copy(io.Discard, reader)
	if ctx.Err() != nil {
		return
	}
	if err != nil {
		log.Printf("parent pipe read failed; shutting down: %v", err)
	} else {
		log.Print("parent pipe closed; shutting down")
	}
	cancel()
}

func removeUnixSocket(target listenerTarget) error {
	if target.Network != "unix" {
		return nil
	}
	if err := os.Remove(target.Address); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove Unix socket %s: %w", target.Address, err)
	}
	return nil
}

func secureUnixSocket(target listenerTarget) error {
	if target.Network != "unix" {
		return nil
	}
	if err := os.Chmod(target.Address, 0o600); err != nil {
		return fmt.Errorf("secure Unix socket %s: %w", target.Address, err)
	}
	return nil
}
