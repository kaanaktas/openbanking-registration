package client

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGenerateJwt(t *testing.T) {

	var claims = map[string]interface{}{
		"grant_types":                  []string{"authorization_code", "refresh_token", "client_credentials"},
		"redirect_uris":                []string{"redirect_uris_1", "redirect_uris_2"},
		"application_type":             "web",
		"iss":                          "iss",
		"token_endpoint_auth_method":   "tls_client_auth",
		"tls_client_auth_dn":           "tls_client_auth_dn",
		"software_id":                  "software_id",
		"software_statement":           "test_ssa",
		"aud":                          "https://obp-api.danskebank.com/open-banking/private",
		"scope":                        "openid accounts payments",
		"jti":                          "40ec08a9-8645-4e4a-ae90-21c473a2a0b8",
		"id_token_signed_response_alg": "PS256",
		"request_object_signing_alg":   "PS256",
		"iat":                          1582717153,
		"exp":                          1582725153,
	}

	_ = os.Setenv("OB_SIGN_KEY", "./testdata/test_key.pem")
	_ = os.Setenv("KID", "kid_test")

	type args struct {
		claims jwt.MapClaims
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"validate_token", args{claims}, "Token is expired", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, vError := GenerateJwt(tt.args.claims)
			if (vError != nil) != tt.wantErr {
				t.Errorf("generateJwt() error = %v\n, wantErr %v", vError, tt.wantErr)
				return
			}

			token, vError := jwt.Parse(got, func(t *jwt.Token) (interface{}, error) {
				certData, _ := ioutil.ReadFile("./testdata/test_cert.pem")
				cert, err := jwt.ParseRSAPublicKeyFromPEM(certData)
				if err != nil {
					log.Fatalln("couldn't retrieve the pem file.", err)
				}
				return cert, err
			})

			if token != nil && token.Header["kid"] != "kid_test" {
				log.Fatalln("kid value doesn't match.")
			}

			if vError != nil && vError.Error() != tt.want {
				t.Errorf("generateJwt() got = %v\n,want = %v", got, tt.want)
			}
		})
	}
}
