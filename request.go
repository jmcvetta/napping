// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Opts struct {
	Userinfo *url.Userinfo // Optional username/password to authenticate this request
	Params   *Params       // URL parameters for GET requests (ignored otherwise)
	Header   *http.Header  // HTTP Headers to use (will override defaults)

	// OPTIONAL - ExpectedStatus is the HTTP status code we expect the server
	// to return on a successful request.  If ExpectedStatus is non-zero and
	// server returns a different code, Send() will return a BadStatus error.
	ExpectedStatus int
}

// A Request describes an HTTP request to be executed, data structures into
// which the result will be unmarshalled, and the server's response. By using
// a  single object for both the request and the response we allow easy access
// to Result and Error objects without needing type assertions.
type Request struct {
	Opts
	Url    string      // Raw URL string
	Method string      // HTTP method to use
	Data   interface{} // Data to JSON-encode and POST

	// Result is a pointer to a data structure.  On success, response from
	// server is unmarshalled into Result.
	Result interface{}

	// The following fields are populated by Send().
	timestamp time.Time      // Time when HTTP request was sent
	rawText   string         // Body of server's response (JSON or otherwise)
	status    int            // HTTP status for executed request
	response  *http.Response // Response object from http package
}

// Timestamp returns the time when HTTP request was sent.
func (r *Request) Timestamp() time.Time {
	return r.timestamp
}

// RawText returns the body of the server's response as raw text.
func (r *Request) RawText() string {
	return r.rawText
}

// Status returns the HTTP status for the executed request, or 0 if request has
// not yet been sent.
func (r *Request) Status() int {
	return r.status
}

// HttpResponse returns the underlying Response object from http package.
func (r *Request) HttpResponse() *http.Response {
	return r.response
}

// Unmarshal parses the JSON-encoded data in the server's response, and stores
// the result in the value pointed to by v.
func (r *Request) Unmarshall(v interface{}) error {
	return json.Unmarshal(r.response.Body, v)
}
