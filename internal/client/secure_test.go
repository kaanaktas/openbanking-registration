package client

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

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

	os.Setenv("CLIENT_CA_CERT_PEM", "./testdata/ob_root_ca.cer")
	os.Setenv("CLIENT_CERT_PEM", "./testdata/BZ5TXmhW1hC6NhqFVB5lURIWzsk.pem")
	os.Setenv("CLIENT_KEY_PEM", "./testdata/zDcwMHjgbPP3ETGEO4tCV9.key")
	client, err := NewSecureHttpClient()
	if err != nil {
		t.Fatalf("could not create secure client, %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := client.CallService(tt.args.endpoint, tt.args.payload); !strings.Contains(got, tt.want) {
				t.Error(err)
				t.Errorf("callService() = %v, want %v", got, tt.want)
			}
		})
	}
}
