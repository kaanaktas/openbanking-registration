package client

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
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

	var key *rsa.PrivateKey
	if value, found := cacheMem.Get("obSignKey"); found {
		key = value.(*rsa.PrivateKey)
	} else {
		signingKey := os.Getenv("OB_SIGN_KEY")
		keyData, _ := ioutil.ReadFile(signingKey)
		var err error
		key, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			return "", fmt.Errorf("couldn't parse the private key. %w", err)
		}
		cacheMem.Set("obSignKey", key, cache.NoExpiration)
	}

	kid := os.Getenv("KID")

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = kid

	return token.SignedString(key)
}
