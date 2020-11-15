package register

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kaanaktas/openbanking-registration/internal/client"
	"github.com/kaanaktas/openbanking-registration/internal/config"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func Register(c echo.Context) error {
	aspspId := c.Param("aspsp")

	register, err := config.LoadConfig(aspspId)
	if err != nil {
		return fmt.Errorf("error while loading aspsp config from config.LoadConfig(). %w", err)
	}

	ssa := c.QueryParam("ssa")
	if ssa == ""{
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":       false,
			"errorMessage": "ssa shouldn't be empty",
		})
	}
	claims := createRegisterPayload(uuid.New().String(), createLongTime(0), createLongTime(60), ssa, register)
	signedPayload, err := client.GenerateJwt(claims)
	if err != nil {
		log.Printf("error while generating jwt, %v", err)

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":       false,
			"errorMessage": err.Error(),
		})
	}

	httpClient, err := client.NewSecureHttpClient()
	if err != nil {
		log.Printf("error while creating new secure http httpClient, %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":       false,
			"errorMessage": err.Error(),
		})
	}

	resp, err := httpClient.CallService(register.Endpoint, []byte(signedPayload))
	if err != nil {
		log.Printf("error while calling backend, %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":       false,
			"errorMessage": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func createRegisterPayload(jti string, iat, exp int64, ssa string, register *config.Register) map[string]interface{} {
	claims := map[string]interface{}{
		"grant_types":                  register.GrantTypes,
		"redirect_uris":                register.RedirectUris,
		"application_type":             register.ApplicationType,
		"iss":                          register.Iss,
		"token_endpoint_auth_method":   register.TokenEndpointAuthMethod,
		"tls_client_auth_subject_dn":   register.TlsClientAuthSubjectDn,
		"software_id":                  register.SoftwareId,
		"software_statement":           ssa,
		"aud":                          register.Aud,
		"scope":                        register.Scope,
		"jti":                          jti,
		"id_token_signed_response_alg": register.IdTokenSignedResponseAlg,
		"request_object_signing_alg":   register.RequestObjectSigningAlg,
		"iat":                          iat,
		"exp":                          exp,
	}

	return claims
}

func createLongTime(addMinute time.Duration) int64 {
	return time.Now().Add(time.Minute * addMinute).Unix()
}
