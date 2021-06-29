package pkg

import "time"

type Config struct {
	Address string        `yaml:"address,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
	TLS     TLSConfig     `yaml:"tls,omitempty"`
	Auth    AuthConfig    `yaml:"auth,omitempty"`
}

type TLSConfig struct {
	CAs        []string `yaml:"cas,omitempty"`
	Cert       string   `yaml:"cert,omitempty"`
	Key        string   `yaml:"key,omitempty"`
	SkipVerify bool     `yaml:"skipVerify,omitempty"`
}

type AuthConfig struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	ClientID string `yaml:"clientId,omitempty"`
}
