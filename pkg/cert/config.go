package cert

type APIVersion string

const (
	APIVersionV1 APIVersion = "v1"
)

type CSRConfig struct {
	APIVersion   APIVersion `json:"apiVersion"`
	CommonName   string     `json:"commonName"`
	Organization string     `json:"organization"`
	Country      string     `json:"country"`
	Province     string     `json:"province"`
	Locality     string     `json:"locality"`
}
