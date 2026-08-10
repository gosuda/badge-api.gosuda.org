package envaddr

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		port        *string
		host        *string
		ip          *string
		defaultAddr string
		want        string
	}{
		{name: "default", defaultAddr: ":8443", want: ":8443"},
		{name: "port", port: new("8080"), defaultAddr: ":8443", want: ":8080"},
		{name: "host and port", port: new("8080"), host: new("localhost"), defaultAddr: ":8443", want: "localhost:8080"},
		{name: "IPv4 and port", port: new("8080"), ip: new("127.255.255.255"), defaultAddr: ":8443", want: "127.255.255.255:8080"},
		{name: "IPv6 and port", port: new("8080"), ip: new("::1"), defaultAddr: ":8443", want: "[::1]:8080"},
		{name: "IP overrides host", port: new("8080"), host: new("localhost"), ip: new("127.0.0.1"), defaultAddr: ":8443", want: "127.0.0.1:8080"},
		{name: "invalid IP preserves host", port: new("8080"), host: new("localhost"), ip: new("invalid"), defaultAddr: ":8443", want: "localhost:8080"},
		{name: "host ignored without port", host: new("localhost"), defaultAddr: ":8443", want: ":8443"},
		{name: "IP ignored without port", ip: new("127.0.0.1"), defaultAddr: ":8443", want: ":8443"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setOptionalEnv(t, "PORT", test.port)
			setOptionalEnv(t, "HOST", test.host)
			setOptionalEnv(t, "IP", test.ip)
			if got := Get(test.defaultAddr); got != test.want {
				t.Fatalf("Get(%q) = %q, want %q", test.defaultAddr, got, test.want)
			}
		})
	}
}

func setOptionalEnv(t *testing.T, key string, input *string) {
	t.Helper()
	previous, existed := os.LookupEnv(key)
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, previous)
			return
		}
		_ = os.Unsetenv(key)
	})
	if input == nil {
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		return
	}
	if err := os.Setenv(key, *input); err != nil {
		t.Fatal(err)
	}
}
