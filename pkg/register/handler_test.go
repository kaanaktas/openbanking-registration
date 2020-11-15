package register

import (
	"github.com/kaanaktas/openbanking-registration/internal/config"
	"reflect"
	"testing"
)

var ssa = "ssa_payload"
var jti = "40ec08a9-8645-4e4a-ae90-21c473a2a0b8"
var iat = int64(1582717153)
var exp = int64(1582725153)
var register = config.Register{
	GrantTypes:               []string{"authorization_code", "refresh_token", "client_credentials"},
	RedirectUris:             []string{"https://openbanking.skyplainsoft.com/callback", "https://localhost:8080/callback"},
	ApplicationType:          "web",
	Iss:                      "AdYBJedrsmAJRzS8SsQg2v",
	TokenEndpointAuthMethod:  "tls_client_auth",
	TlsClientAuthSubjectDn:   "CN=AdYBJedrsmAJRzS8SsQg2v, OU=0014H00001lFE4pQAG, O=OpenBanking, C=GB",
	SoftwareId:               "AdYBJedrsmAJRzS8SsQg2v",
	Aud:                      "https://obp-api.danskebank.com/open-banking/private",
	Scope:                    "openid accounts payments",
	IdTokenSignedResponseAlg: "PS256",
	RequestObjectSigningAlg:  "PS256",
}

var claims = map[string]interface{}{
	"grant_types":                  []string{"authorization_code", "refresh_token", "client_credentials"},
	"redirect_uris":                []string{"https://openbanking.skyplainsoft.com/callback", "https://localhost:8080/callback"},
	"application_type":             "web",
	"iss":                          "AdYBJedrsmAJRzS8SsQg2v",
	"token_endpoint_auth_method":   "tls_client_auth",
	"tls_client_auth_subject_dn":   "CN=AdYBJedrsmAJRzS8SsQg2v, OU=0014H00001lFE4pQAG, O=OpenBanking, C=GB",
	"software_id":                  "AdYBJedrsmAJRzS8SsQg2v",
	"software_statement":           ssa,
	"aud":                          "https://obp-api.danskebank.com/open-banking/private",
	"scope":                        "openid accounts payments",
	"jti":                          jti,
	"id_token_signed_response_alg": "PS256",
	"request_object_signing_alg":   "PS256",
	"iat":                          iat,
	"exp":                          exp,
}

func Test_createRegisterPayload(t *testing.T) {
	type args struct {
		jti      string
		iat      int64
		exp      int64
		ssa      string
		register *config.Register
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			"compare_with_ssa",
			args{jti: jti, iat: iat, exp: exp, ssa: ssa, register: &register},
			claims,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createRegisterPayload(tt.args.jti, tt.args.iat, tt.args.exp, tt.args.ssa, tt.args.register); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRegisterPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
