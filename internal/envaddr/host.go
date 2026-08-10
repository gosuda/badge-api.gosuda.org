package envaddr

import (
	"net"
	"os"
)

// Get returns a listener address derived from PORT, HOST, and IP.
// If PORT is not set, Get returns defaultAddr unchanged.
func Get(defaultAddr string) string {
	addr := defaultAddr
	port := ""
	portFound := false
	host := ""

	portVar, exists := os.LookupEnv("PORT")
	if exists {
		port = portVar
		portFound = true
	}
	hostVar, exists := os.LookupEnv("HOST")
	if exists {
		host = hostVar
	}
	ipEnv, exists := os.LookupEnv("IP")
	if exists {
		ip := net.ParseIP(ipEnv)
		if ip != nil {
			if ip.To4() != nil {
				host = ip.String()
			} else if ip.To16() != nil {
				host = "[" + ip.String() + "]"
			}
		}
	}

	if portFound {
		addr = host + ":" + port
	}

	return addr
}
