// Package certificates owns the trust roots shipped with GURL. A private pool
// keeps TLS behavior independent from obsolete or missing stores on old hosts.
package certificates

import (
	"crypto/x509"
	_ "embed"
	"fmt"
)

// Mozilla roots converted and published by https://curl.se/docs/caextract.html.
// The embedded file is MPL 2.0 licensed and remains unmodified.
//
//go:embed cacert.pem
var certificateBundle []byte

// Pool returns a new certificate pool so each request owns its TLS state.
func Pool() (*x509.CertPool, error) {
	return poolFromPEM(certificateBundle)
}

func poolFromPEM(certificatePEM []byte) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certificatePEM) {
		return nil, fmt.Errorf("embedded CA bundle contains no certificates")
	}
	return pool, nil
}
