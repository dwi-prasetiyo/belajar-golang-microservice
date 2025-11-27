package util

import (
	"crypto/tls"
	"crypto/x509"
	"user-service/src/common/log"
)

func CreateServerTlsConf(caCert []byte, serverCert tls.Certificate) *tls.Config {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Logger.Fatal("failed to add CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
}
