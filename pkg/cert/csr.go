package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type CSR struct {
	PrivateKey []byte
	Request    []byte
}

func CreateCSR() (CSR, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return CSR{}, err
	}

	reqTmpl := x509.CertificateRequest{}
	req, err := x509.CreateCertificateRequest(rand.Reader, &reqTmpl, privKey)
	if err != nil {
		return CSR{}, fmt.Errorf("failed to create certificate request: %w", err)
	}

	privKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}); err != nil {
		return CSR{}, fmt.Errorf("failed to encode private key: %w", err)
	}

	reqPEM := new(bytes.Buffer)
	if err := pem.Encode(reqPEM, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: req,
	}); err != nil {
		return CSR{}, fmt.Errorf("failed to encode CSR: %w", err)
	}

	return CSR{
		PrivateKey: privKeyPEM.Bytes(),
		Request:    reqPEM.Bytes(),
	}, nil
}
