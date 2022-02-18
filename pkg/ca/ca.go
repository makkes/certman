package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
)

type CA struct {
	Certificate []byte
	PrivateKey  []byte
}

func CreateCA(cfg *CAConfig) (CA, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: []string{cfg.Organization},
			Country:      []string{cfg.Country},
			Province:     []string{cfg.Province},
			Locality:     []string{cfg.Locality},
		},
		NotBefore:             cfg.NotBefore,
		NotAfter:              cfg.NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return CA{}, fmt.Errorf("failed to create private key: %w", err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return CA{}, fmt.Errorf("failed to create certificate: %w", err)
	}

	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return CA{}, fmt.Errorf("failed to encode certificate: %w", err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}); err != nil {
		return CA{}, fmt.Errorf("failed to encode private key: %w", err)
	}

	return CA{
		Certificate: caPEM.Bytes(),
		PrivateKey:  caPrivKeyPEM.Bytes(),
	}, nil
}
