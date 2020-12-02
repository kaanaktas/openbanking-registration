package client

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func Test_callService(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"name": "test",
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
		{"http_valid_post", args{"https://ob19-rs1.o3bank.co.uk:4501/dynamic-client-registration/v3.1/register", postBody}, "invalid_client"},
	}

	os.Setenv("CLIENT_CA_CERT_PEM", "./testdata/ob_issuer.cer")
	os.Setenv("CLIENT_CERT_PEM", "./testdata/test_cert.pem")
	os.Setenv("CLIENT_KEY_PEM", "./testdata/test_key.pem")

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
