// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"net/http"
	"net/url"
	"time"
)

type Opts struct {
	Userinfo *url.Userinfo // Optional username/password to authenticate this request
	Params   *Params       // URL parameters for GET requests (ignored otherwise)
	Header   *http.Header  // HTTP Headers to use (will override defaults)

	// Pointer to data structure.  On error, response from server is
	// unmarshaled into this structure.
	Error interface{}

	// OPTIONAL - ExpectedStatus is the HTTP status code we expect the server
	// to return on a successful request.  If ExpectedStatus is non-zero and
	// server returns a different code, Send() will return a BadStatus error.
	ExpectedStatus int
}

// A Request describes an HTTP request to be executed, data structures into
// which results and errors will be unmarshalled, and the server's response.
// By using a single object for both the request and the response we allow easy
// access to Result and Error objects without needing type assertions.
type Request struct {
	Opts
	Url    string      // Raw URL string
	Method string      // HTTP method to use
	Data   interface{} // Data to JSON-encode and POST

	// Result is a pointer to a data structure.  On success, response from
	// server is unmarshalled into Result.
	Result interface{}

	// The following fields are populated by Send().
	Timestamp    time.Time      // Time when HTTP request was sent
	RawText      string         // Raw text of server response (JSON or otherwise)
	Status       int            // HTTP status for executed request
	HttpResponse *http.Response // Response object from http package
}
