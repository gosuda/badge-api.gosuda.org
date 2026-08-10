package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestConfiguredListenerTarget(t *testing.T) {
	tests := []struct {
		name     string
		httpAddr *string
		port     *string
		host     *string
		ip       *string
		want     listenerTarget
		wantErr  bool
	}{
		{name: "default TCP", want: listenerTarget{Network: "tcp", Address: ":8080"}},
		{name: "PaaS TCP", port: new("9000"), host: new("localhost"), want: listenerTarget{Network: "tcp", Address: "localhost:9000"}},
		{name: "Unix socket", httpAddr: new("unix:/tmp/badge-api.sock"), port: new("9000"), want: listenerTarget{Network: "unix", Address: "/tmp/badge-api.sock"}},
		{name: "explicit TCP", httpAddr: new("tcp:127.0.0.1:7000"), want: listenerTarget{Network: "tcp", Address: "127.0.0.1:7000"}},
		{name: "plain TCP", httpAddr: new("127.0.0.1:7100"), want: listenerTarget{Network: "tcp", Address: "127.0.0.1:7100"}},
		{name: "empty Unix socket", httpAddr: new("unix:"), wantErr: true},
		{name: "empty explicit TCP", httpAddr: new("tcp:"), wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setOptionalEnvironment(t, httpAddressEnvironment, test.httpAddr)
			setOptionalEnvironment(t, "PORT", test.port)
			setOptionalEnvironment(t, "HOST", test.host)
			setOptionalEnvironment(t, "IP", test.ip)
			got, err := configuredListenerTarget()
			if test.wantErr {
				if err == nil {
					t.Fatalf("configuredListenerTarget() = %#v, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("configuredListenerTarget() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestStdinEOFShutdownEnabled(t *testing.T) {
	setOptionalEnvironment(t, exitOnStdinEOFEnvironment, nil)
	enabled, err := stdinEOFShutdownEnabled()
	if err != nil || enabled {
		t.Fatalf("stdinEOFShutdownEnabled() = %v, %v; want false, nil", enabled, err)
	}

	setOptionalEnvironment(t, exitOnStdinEOFEnvironment, new("true"))
	enabled, err = stdinEOFShutdownEnabled()
	if err != nil || !enabled {
		t.Fatalf("stdinEOFShutdownEnabled() = %v, %v; want true, nil", enabled, err)
	}

	setOptionalEnvironment(t, exitOnStdinEOFEnvironment, new("invalid"))
	if _, err = stdinEOFShutdownEnabled(); err == nil {
		t.Fatal("stdinEOFShutdownEnabled() accepted an invalid boolean")
	}
}

func TestShutdownOnStdinEOF(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownOnStdinEOF(ctx, strings.NewReader(""), cancel)
	select {
	case <-ctx.Done():
	default:
		t.Fatal("stdin EOF did not cancel the application context")
	}
}

func TestRemoveUnixSocket(t *testing.T) {
	path := t.TempDir() + "/badge.sock"
	if err := os.WriteFile(path, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := removeUnixSocket(listenerTarget{Network: "unix", Address: path}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("socket path still exists: %v", err)
	}
}

func setOptionalEnvironment(t *testing.T, key string, value *string) {
	t.Helper()
	previous, existed := os.LookupEnv(key)
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, previous)
			return
		}
		_ = os.Unsetenv(key)
	})
	if value == nil {
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		return
	}
	if err := os.Setenv(key, *value); err != nil {
		t.Fatal(err)
	}
}
