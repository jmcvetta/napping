// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package restclient provides a simple client library for interacting with
// RESTful APIs.
package restclient

import (
	// "log"
	"bytes"
	// "reflect"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Method string

var (
	GET    = Method("GET")
	PUT    = Method("PUT")
	POST   = Method("POST")
	DELETE = Method("DELETE")
)

type RestRequest struct {
	Url     string            // Raw URL string
	Method  Method            // HTTP method to use 
	Params  map[string]string // URL parameters for GET requests (ignored otherwise)
	Headers *http.Header      // HTTP Headers to use (will override defaults)
	Data    interface{}       // Data to JSON-encode and include with call
	Result  interface{}       // JSON-encoded data in respose will be unmarshalled into Result
	Error   interface{}       // If server returns error status, JSON-encoded response data will be unmarshalled into Error
	RawText string            // Gets populated with raw text of server response
}

// Client is a REST client.
type Client struct {
	HttpClient *http.Client
}

// New returns a new Client instance.
func New() *Client {
	return &Client{
		HttpClient: new(http.Client),
	}
}

// Do executes an HTTP REST request
func (c *Client) Do(r *RestRequest) (status int, err error) {
	//
	// Create a URL object from the raw url string.  This will allow us to compose
	// query parameters programmatically and be guaranteed of a well-formed URL.
	//
	u, err := url.Parse(r.Url)
	if err != nil {
		return
	}
	//
	// If we are making a GET request and the user populated the Params field, then
	// add the params to the URL's querystring.
	//
	if r.Method == GET && r.Params != nil {
		vals := u.Query()
		for k, v := range r.Params {
			vals.Set(k, v)
		}
		u.RawQuery = vals.Encode()
	}
	//
	// Create a Request object; if populated, Data field is JSON encoded as request
	// body
	//
	m := string(r.Method)
	var req *http.Request
	if r.Data == nil {
		req, err = http.NewRequest(m, u.String(), nil)
	} else {
		var b []byte
		b, err = json.Marshal(r.Data)
		if err != nil { // Create a URL object from the raw url string.  This will allow us to compose URL parameters programmatically and be
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(m, u.String(), buf)
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return
	}
	//
	// If Accept header is unset, set it for JSON.
	//
	if req.Header.Get("Accept") == "" {
		req.Header.Add("Accept", "application/json")
	}
	//
	// Execute the HTTP request
	//
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return
	}
	status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	r.RawText = string(data)
	if status >= 200 && status < 300 {
		err = json.Unmarshal(data, &r.Result)
	} else {
		err = json.Unmarshal(data, &r.Error)
	}
	return
}
