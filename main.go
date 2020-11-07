package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var headers = map[string]interface{}{
	"typ": "JWT",
	"alg": "PS256",
	"kid": "KID",
}

var endpoints = map[string]string{
	"danske": "https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v1.0/thirdparty/register",
	"hsbc":   "hsbc",
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

func doRegister(c echo.Context) error {
	endpoint := endpoints[c.Param("aspsp")]
	ssa := c.QueryParam("ssa")
	claims := createRegisterPayload(ssa)
	signedPayload, err := generateJwt(claims)
	if err != nil {
		log.Panic(err)
		return err
	}

	resp, err := callService(endpoint, []byte(signedPayload))
	if err != nil {
		log.Panic(err)
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func createRegisterPayload(ssa string) map[string]interface{} {
	claims := map[string]interface{}{
		"grant_types":                  []string{"authorization_code", "refresh_token", "client_credentials"},
		"redirect_uris":                []string{"https://openbanking.skyplainsoft.com/callback", "https://localhost:8080/callback"},
		"application_type":             "web",
		"iss":                          "AdYBJedrsmAJRzS8SsQg2v",
		"token_endpoint_auth_method":   "tls_client_auth",
		"tls_client_auth_dn":           "CN=AdYBJedrsmAJRzS8SsQg2v, OU=0014H00001lFE4pQAG, O=OpenBanking, C=GB",
		"software_id":                  "AdYBJedrsmAJRzS8SsQg2v",
		"software_statement":           ssa,
		"aud":                          "https://obp-api.danskebank.com/open-banking/private",
		"scope":                        "openid accounts payments",
		"jti":                          "40ec08a9-8645-4e4a-ae90-21c473a2a0b8",
		"id_token_signed_response_alg": "PS256",
		"request_object_signing_alg":   "PS256",
		"iat":                          createLongTime(0),
		"exp":                          createLongTime(60),
	}

	return claims
}

func generateJwt(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header = headers

	keyData, _ := ioutil.ReadFile("certs/ob-transport/signing.key")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		fmt.Println(err)
	}

	return token.SignedString(key)
}

func callService(endpoint string, payload []byte) (string, error) {
	caCert, err := ioutil.ReadFile("certs/ob-transport/ob_root_ca.cer")
	if err != nil {
		return "", err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair("certs/ob-transport/BZ5TXmhW1hC6NhqFVB5lURIWzsk.pem", "certs/ob-transport/zDcwMHjgbPP3ETGEO4tCV9.key")
	tlsConfig := tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{clientCert},
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := http.Client{
		Transport: &transport,
	}

	responseBody := bytes.NewBuffer(payload)
	resp, err := client.Post(endpoint, "application/jwt", responseBody)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func createLongTime(addMinute time.Duration) int64 {
	return time.Now().Add(time.Minute * addMinute).Unix()
}
