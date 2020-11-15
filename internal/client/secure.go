package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SecureClient struct {
	*http.Client
}

var client *SecureClient

func NewSecureHttpClient() (*SecureClient, error) {
	if client != nil {
		return client, nil
	}

	serverCACertPEM := os.Getenv("CLIENT_CA_CERT_PEM")
	clientCertPEM := os.Getenv("CLIENT_CERT_PEM")
	clientKeyPEM := os.Getenv("CLIENT_KEY_PEM")

	caCert, err := ioutil.ReadFile(serverCACertPEM)
	if err != nil {
		return nil, fmt.Errorf("read cert file, %w", err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("load x509 certificate, %w", err)
	}

	transport := http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCert},
		},
	}

	client = &SecureClient{
		&http.Client{
			Transport: &transport,
		},
	}

	return client, nil
}

const MIMEApplicationJWT = "application/jwt"

func (s *SecureClient) CallService(endpoint string, payload []byte) (string, error) {
	responseBody := bytes.NewBuffer(payload)

	resp, err := s.Client.Post(endpoint, MIMEApplicationJWT, responseBody)

	if err != nil {
		return "", fmt.Errorf("error new secure post, %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error read response body,%w", err)
	}

	return string(body), nil
}
