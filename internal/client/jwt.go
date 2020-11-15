package client

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
)

func GenerateJwt(claims jwt.MapClaims) (string, error) {
	signingAlgorithm := claims["request_object_signing_alg"]
	var signingMethod jwt.SigningMethod

	if signingAlgorithm == "PS256" {
		//below 2 lines fix the wrong signature problem for PS-*.
		//follow up issue: https://github.com/dgrijalva/jwt-go/pull/305
		signingMethodPS256 := jwt.SigningMethodPS256
		signingMethodPS256.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
		signingMethod = signingMethodPS256
	} else if signingAlgorithm == "RS256" {
		signingMethod = jwt.SigningMethodRS256
	} else {
		return "", fmt.Errorf("unsupported signing algorithm %v", signingAlgorithm)
	}

	signingKey := os.Getenv("OB_SIGN_KEY")
	kid := os.Getenv("KID")

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = kid

	keyData, _ := ioutil.ReadFile(signingKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)

	if err != nil {
		return "", fmt.Errorf("couldn't parse the private key. %w", err)
	}

	return token.SignedString(key)
}
