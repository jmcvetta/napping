// Package httpclient provides a way to get hold of an http.Client instance that
// will timeout requests. Otherwise, they will not time out.  This code snippet
// was based on https://gist.github.com/dmichael/5710968
package httpclient

import (
	"net"
	"net/http"
	"time"
)

// Create an http Client that times out in two ways:
// 1) connectTimeout - times out trying to connect
// 2) readWriteTimeout - times out reading or writing data
// If either is set to 0, the http client will not timeout for its given reason.
func NewWithTimeout(connectTimeout, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{Transport: TimeoutTransport(connectTimeout, readWriteTimeout)}
}

func TimeoutTransport(connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	newTimeoutConnection := func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		if (readWriteTimeout > 0) {
			conn.SetDeadline(time.Now().Add(readWriteTimeout))
		}
		return conn, nil
	}
	return &http.Transport{Dial: newTimeoutConnection}
}
