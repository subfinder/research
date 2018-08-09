package core

import (
	"net"
	"net/http"
	"time"
)

// HTTPClient is a reusable component that can be used in sources.
var HTTPClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   30 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	},
}
