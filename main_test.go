package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"strings"
	"testing"
)

const validToken = "eyJhbGciOiJIUzI1NiIsImtpZCI6IktJRCIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbl90eXBlIjoid2ViIiwiYXVkIjoiaHR0cHM6Ly9vYnAtYXBpLmRhbnNrZWJhbmsuY29tL29wZW4tYmFua2luZy9wcml2YXRlIiwiZXhwIjoxNTgyNzI1MTUzLCJncmFudF90eXBlcyI6WyJhdXRob3JpemF0aW9uX2NvZGUiLCJyZWZyZXNoX3Rva2VuIiwiY2xpZW50X2NyZWRlbnRpYWxzIl0sImlhdCI6MTU4MjcxNzE1MywiaWRfdG9rZW5fc2lnbmVkX3Jlc3BvbnNlX2FsZyI6IlBTMjU2IiwiaXNzIjoiQWRZQkplZHJzbUFKUnpTOFNzUWcydiIsImp0aSI6IjQwZWMwOGE5LTg2NDUtNGU0YS1hZTkwLTIxYzQ3M2EyYTBiOCIsInJlZGlyZWN0X3VyaXMiOlsiaHR0cHM6Ly9vcGVuYmFua2luZy5za3lwbGFpbnNvZnQuY29tL2NhbGxiYWNrIiwiaHR0cHM6Ly9sb2NhbGhvc3Q6ODA4MC9jYWxsYmFjayJdLCJyZXF1ZXN0X29iamVjdF9zaWduaW5nX2FsZyI6IlBTMjU2Iiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMgcGF5bWVudHMiLCJzb2Z0d2FyZV9pZCI6IkFkWUJKZWRyc21BSlJ6UzhTc1FnMnYiLCJzb2Z0d2FyZV9zdGF0ZW1lbnQiOiJ0ZXN0X3NzYSIsInRsc19jbGllbnRfYXV0aF9kbiI6IkNOPUFkWUJKZWRyc21BSlJ6UzhTc1FnMnYsIE9VPTAwMTRIMDAwMDFsRkU0cFFBRywgTz1PcGVuQmFua2luZywgQz1HQiIsInRva2VuX2VuZHBvaW50X2F1dGhfbWV0aG9kIjoidGxzX2NsaWVudF9hdXRoIn0.saQDpBIx_IE6oOrGE6gdQJrePOIZa6ESOOvZDZl8Xok"

var claims = map[string]interface{}{
	"grant_types":                  []string{"authorization_code", "refresh_token", "client_credentials"},
	"redirect_uris":                []string{"https://openbanking.skyplainsoft.com/callback", "https://localhost:8080/callback"},
	"application_type":             "web",
	"iss":                          "AdYBJedrsmAJRzS8SsQg2v",
	"token_endpoint_auth_method":   "tls_client_auth",
	"tls_client_auth_dn":           "CN=AdYBJedrsmAJRzS8SsQg2v, OU=0014H00001lFE4pQAG, O=OpenBanking, C=GB",
	"software_id":                  "AdYBJedrsmAJRzS8SsQg2v",
	"software_statement":           "test_ssa",
	"aud":                          "https://obp-api.danskebank.com/open-banking/private",
	"scope":                        "openid accounts payments",
	"jti":                          "40ec08a9-8645-4e4a-ae90-21c473a2a0b8",
	"id_token_signed_response_alg": "PS256",
	"request_object_signing_alg":   "PS256",
	"iat":                          1582717153,
	"exp":                          1582725153,
}

func Test_createRegisterPayload(t *testing.T) {
	type args struct {
		ssa string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			"compare_with_payload",
			args{ssa: "test_ssa"},
			claims,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createRegisterPayload(tt.args.ssa); !reflect.DeepEqual(got["ssa"], tt.want["ssa"]) {
				t.Errorf("createRegisterPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateJwt(t *testing.T) {
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
			"compare_with_token", args{claims}, validToken, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJwt(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJwt() error = %v\n, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateJwt() got = %v\n,want = %v", got, tt.want)
			}
		})
	}
}

func Test_callService(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"name":  "test",
		"email": "test@test.com",
	})

	type args struct {
		endpoint string
		payload  []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"http_valid_post", args{"https://sandbox-obp-api.danskebank.com/sandbox-open-banking/v1.0/thirdparty/register", postBody}, "{"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := callService(tt.args.endpoint, tt.args.payload); !strings.Contains(got, tt.want) {
				t.Error(err)
				t.Errorf("callService() = %v, want %v", got, tt.want)
			}
		})
	}
}