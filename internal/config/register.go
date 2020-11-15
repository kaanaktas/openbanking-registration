package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type config struct {
	Register Register `yaml:"register"`
}

type Register struct {
	Endpoint                 string   `yaml:"endpoint"`
	GrantTypes               []string `yaml:"grantTypes"`
	RedirectUris             []string `yaml:"redirectUris"`
	ApplicationType          string   `yaml:"applicationType"`
	Iss                      string   `yaml:"iss"`
	TokenEndpointAuthMethod  string   `yaml:"tokenEndpointAuthMethod"`
	TlsClientAuthSubjectDn   string   `yaml:"tlsClientAuthSubjectDn"`
	SoftwareId               string   `yaml:"softwareId"`
	Aud                      string   `yaml:"aud"`
	Scope                    string   `yaml:"scope"`
	IdTokenSignedResponseAlg string   `yaml:"idTokenSignedResponseAlg"`
	RequestObjectSigningAlg  string   `yaml:"requestObjectSigningAlg"`
}

var cache = map[string]Register{}

func LoadConfig(aspspId string) (*Register, error) {
	if item, ok := cache[aspspId]; ok {
		return &item, nil
	}

	f, err := os.Open("./aspsp/" + aspspId + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("error open yaml file, %w", err)
	}
	defer f.Close()

	var aspsp config
	if err := yaml.NewDecoder(f).Decode(&aspsp); err != nil {
		return nil, fmt.Errorf("error yaml decode, %w", err)
	}
	cache[aspspId] = aspsp.Register

	return &aspsp.Register, nil
}
