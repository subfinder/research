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
			Timeout:   30 * time.Second,
			KeepAlive: 600 * time.Second,
		}).Dial,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		TLSHandshakeTimeout:   30 * time.Second,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	},
}
