package client

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
)

var headers = map[string]interface{}{
	"typ": "JWT",
	"alg": "PS256",
	"kid": "O9JRkzXnFf6AK7H2kq2UI_Gv2I8",
}

func GenerateJwt(claims jwt.MapClaims) (string, error) {
	signingKey := os.Getenv("OB_SIGN_KEY")

	//Below fixes the wrong signature problem for PS-*.
	//this is the follow up issue -> https://github.com/dgrijalva/jwt-go/pull/305
	signingMethodPS256 := jwt.SigningMethodPS256
	signingMethodPS256.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header = headers

	keyData, _ := ioutil.ReadFile(signingKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}
