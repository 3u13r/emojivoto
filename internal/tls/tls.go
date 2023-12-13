package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func LoadTLSConfig() *tls.Config {
	certPath, ok := os.LookupEnv("EDG_CERT_PATH")
	if !ok {
		panic("EDG_CERT_PATH environment variable must be set")
	}
	keyPath, ok := os.LookupEnv("EDG_KEY_PATH")
	if !ok {
		panic("EDG_KEY_PATH environment variable must be set")
	}
	caCertPath, ok := os.LookupEnv("EDG_CA_PATH")
	if !ok {
		panic("EDG_CA_PATH environment variable must be set")
	}
	disableClientAuth := os.Getenv("EDG_DISABLE_CLIENT_AUTH")

	// Load TLS credentials
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	clientAuth := tls.RequireAndVerifyClientCert
	if disableClientAuth == "true" {
		clientAuth = tls.NoClientCert
	}

	// Load CA cert
	rootCAs := x509.NewCertPool()
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		panic("failed to append CA cert")
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   clientAuth,
		RootCAs:      rootCAs,
		ClientCAs:    rootCAs,
	}

	return config
}
