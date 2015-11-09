// Provides basic utilities to get server crypto config.
package crypto

import (
	"crypto/tls"
	"os"
)

func IsEnabled() bool {
	filename := GetX509KeyPair()[0]
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

func GetX509KeyPair() [2]string {
	return [2]string{"crypto/server.pem", "crypto/server.key"}
}

func GetServerTlsConfig() (*tls.Config, error) {
	keypair := GetX509KeyPair()
	cert, err := tls.LoadX509KeyPair(keypair[0], keypair[1])
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func GetClientTlsConfig() (*tls.Config, error) {
	if !IsEnabled() {
		return nil, nil
	}
	return &tls.Config{InsecureSkipVerify: true}, nil
}
