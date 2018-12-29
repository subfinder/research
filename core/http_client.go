package core

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/viki-org/dnscache"
)

var dnsCache = dnscache.New(5 * time.Minute)

// HTTPClient is a reusable component that can be used in sources.
var HTTPClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 64,
		Dial: func(network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ip, err := dnsCache.FetchOneString(address[:separator])
			if err != nil {
				return nil, err
			}
			return net.Dial("tcp", ip+address[separator:])
		},
		//Dial: (&net.Dialer{
		//	Timeout:   15 * time.Second,
		//	KeepAlive: 60 * time.Second,
		//}).Dial,
		MaxIdleConns: 100,
		//MaxIdleConnsPerHost: 2,
		//MaxConnsPerHost:       1,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// var HTTPClient = &http.Client{
//	Transport: &http.Transport{
//		Dial: (&net.Dialer{
//			Timeout:   30 * time.Second,
//			KeepAlive: 600 * time.Second,
//		}).Dial,
//		MaxIdleConns:          100,
//		MaxIdleConnsPerHost:   100,
//		TLSHandshakeTimeout:   30 * time.Second,
//		IdleConnTimeout:       30 * time.Second,
//		ResponseHeaderTimeout: 60 * time.Second,
//		ExpectContinueTimeout: 30 * time.Second,
//	},
//}
