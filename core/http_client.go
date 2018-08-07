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
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}
