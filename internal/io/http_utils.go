package io

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

// Download downloads content from the given URL and returns it as an io.ReadCloser.
// It uses an insecure TLS configuration to work around certificate issues (e.g., thapar.edu).
// WARNING: This is insecure and only used as a workaround for certificate issues.
func Download(url string) (io.ReadCloser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	return resp.Body, nil
}
