package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

func NewTLSConfig(c TLSConfig) (*tls.Config, error) {
	if c.Key == "" && c.Cert == "" {
		log.Debug("tls key/cert not configured, will not create TLS config")
		return nil, nil
	}

	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: c.SkipVerify,
	}

	if len(c.CAs) > 0 {
		log.Debug("creating CA cert pool for TLS config")
		pool, err := makeCertPool(c.CAs...)
		if err != nil {
			return nil, err
		}
		cfg.RootCAs = pool
	}

	if c.Cert != "" && c.Key != "" {
		if err := loadCertificate(cfg, c.Cert, c.Key); err != nil {
			return nil, err
		}
	} else {
		log.WithFields(log.Fields{
			"key":  c.Key,
			"cert": c.Cert,
		}).Error("both cert and key are required")
		return nil, errors.New("both cert and key are required for tls config")
	}

	return cfg, nil
}

func makeCertPool(certs ...string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, c := range certs {
		pem, err := ioutil.ReadFile(c)
		if err != nil {
			return nil, errors.WithMessagef(err, "error reading root CA from file '%s'", c)
		}
		ok := pool.AppendCertsFromPEM(pem)
		if !ok {
			log.WithField("ca", c).Error("[openconfig] failed to add CA PEM to cert pool")
			return nil, errors.New("failed to append cert from PEM")
		}
	}
	return pool, nil
}

func loadCertificate(c *tls.Config, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"cert":  certFile,
			"key":   keyFile,
		}).Error("[openconfig] failed to load client certificate")
		return errors.WithMessage(err, "failed to load client certificate")
	}
	c.Certificates = append(c.Certificates, cert)
	return nil
}
