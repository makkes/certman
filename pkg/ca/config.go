package ca

import "time"

type APIVersion string

const (
	APIVersionV1 APIVersion = "v1"
)

type CAConfig struct {
	APIVersion   APIVersion `json:"apiVersion"`
	CommonName   string     `json:"commonName"`
	Organization string     `json:"organization"`
	Country      string     `json:"country"`
	Province     string     `json:"province"`
	Locality     string     `json:"locality"`
	NotBefore    time.Time  `json:"notBefore"`
	NotAfter     time.Time  `json:"notAfter"`
}
