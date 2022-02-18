package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"time"
)

func CreateCertificate(csrPEM, caPEM, caKeyPEM []byte) ([]byte, error) {
	pemCert, _ := pem.Decode(csrPEM)
	if pemCert == nil {
		return nil, fmt.Errorf("no PEM data found in CSR input")
	}
	csr, err := x509.ParseCertificateRequest(pemCert.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	pemCA, _ := pem.Decode(caPEM)
	if pemCA == nil {
		return nil, fmt.Errorf("no PEM data found in CA input")
	}
	caCert, err := x509.ParseCertificate(pemCA.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	pemCAKey, _ := pem.Decode(caKeyPEM)
	if pemCAKey == nil {
		return nil, fmt.Errorf("no PEM data found in CA key input")
	}
	caPrivKey, err := x509.ParsePKCS1PrivateKey(pemCAKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA private key: %w", err)
	}

	ser, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	cert := x509.Certificate{
		SerialNumber: ser,
		Issuer:       caCert.Subject,
		Subject:      csr.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 90),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		IsCA:     false,
		DNSNames: []string{csr.Subject.CommonName},
	}

	signedCert, err := x509.CreateCertificate(rand.Reader, &cert, caCert, csr.PublicKey, caPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate from CSR: %w", err)
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signedCert,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode certificate: %w", err)
	}

	return certPEM.Bytes(), nil
}
