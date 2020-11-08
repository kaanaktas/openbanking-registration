package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	MIMEApplicationJWT       = "application/jwt"
	ClientKeyPEM             = "certs/ob-transport/zDcwMHjgbPP3ETGEO4tCV9.key"
	ClientCertPEM            = "certs/ob-transport/BZ5TXmhW1hC6NhqFVB5lURIWzsk.pem"
	ServerCACertPEM          = "certs/ob-transport/ob_root_ca.cer"
	OpenBankingSigningKeyPEM = "certs/ob-transport/signing.key"
)

var headers = map[string]interface{}{
	"typ": "JWT",
	"alg": "PS256",
	"kid": "O9JRkzXnFf6AK7H2kq2UI_Gv2I8",
}

type Register struct {
	Endpoint                 string   `yaml:"endpoint"`
	GrantTypes               []string `yaml:"grantTypes"`
	RedirectUris             []string `yaml:"redirectUris"`
	ApplicationType          string   `yaml:"applicationType"`
	Iss                      string   `yaml:"iss"`
	TokenEndpointAuthMethod  string   `yaml:"tokenEndpointAuthMethod"`
	TlsClientAuthDn          string   `yaml:"tlsClientAuthDn"`
	SoftwareId               string   `yaml:"softwareId"`
	Aud                      string   `yaml:"aud"`
	Scope                    string   `yaml:"scope"`
	IdTokenSignedResponseAlg string   `yaml:"idTokenSignedResponseAlg"`
	RequestObjectSigningAlg  string   `yaml:"requestObjectSigningAlg"`
}

type Aspsp struct {
	Register Register `yaml:"register"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/:aspsp/register", doRegister)

	log.Println("Listing for requests at http://localhost:8080/")
	e.Logger.Fatal(e.Start(":8080"))
}

func readAspspConfiguration(aspspId string) (Register, error) {
	f, err := os.Open("aspsp/" + aspspId + ".yaml")
	if err != nil {
		return Register{}, err
	}
	defer f.Close()

	var aspsp Aspsp
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&aspsp)
	if err != nil {
		return Register{}, err
	}

	return aspsp.Register, nil
}

func doRegister(c echo.Context) error {
	aspspId := c.Param("aspsp")
	register, err := readAspspConfiguration(aspspId)
	if err != nil {
		return err
	}

	ssa := c.QueryParam("ssa")
	claims := createRegisterPayload(ssa, register)
	signedPayload, err := generateJwt(claims)
	if err != nil {
		return err
	}

	resp, err := callService(register.Endpoint, []byte(signedPayload))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func createRegisterPayload(ssa string, register Register) map[string]interface{} {
	claims := map[string]interface{}{
		"grant_types":                  register.GrantTypes,
		"redirect_uris":                register.RedirectUris,
		"application_type":             register.ApplicationType,
		"iss":                          register.Iss,
		"token_endpoint_auth_method":   register.TokenEndpointAuthMethod,
		"tls_client_auth_dn":           register.TlsClientAuthDn,
		"software_id":                  register.SoftwareId,
		"software_statement":           ssa,
		"aud":                          register.Aud,
		"scope":                        register.Scope,
		"jti":                          uuid.New(),
		"id_token_signed_response_alg": register.IdTokenSignedResponseAlg,
		"request_object_signing_alg":   register.RequestObjectSigningAlg,
		"iat":                          createLongTime(0),
		"exp":                          createLongTime(60),
	}

	return claims
}

func generateJwt(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header = headers

	keyData, _ := ioutil.ReadFile(OpenBankingSigningKeyPEM)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

func callService(endpoint string, payload []byte) (string, error) {
	caCert, err := ioutil.ReadFile(ServerCACertPEM)
	if err != nil {
		return "", err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(ClientCertPEM, ClientKeyPEM)
	tlsConfig := tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{clientCert},
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := http.Client{
		Transport: &transport,
	}

	responseBody := bytes.NewBuffer(payload)
	resp, err := client.Post(endpoint, MIMEApplicationJWT, responseBody)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func createLongTime(addMinute time.Duration) int64 {
	return time.Now().Add(time.Minute * addMinute).Unix()
}
