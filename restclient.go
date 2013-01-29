// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package restclient provides a simple client library for interacting with
// RESTful APIs.
package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
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
	Status  int               // HTTP status for executed request
}

// Client is a REST client.
type Client struct {
	HttpClient   *http.Client
	DefaultError interface{}
}

// New returns a new Client instance.
func New() *Client {
	return &Client{
		HttpClient: new(http.Client),
	}
}

// Do executes a REST request.
func (c *Client) Do(r *RestRequest) (status int, err error) {
	if r.Error == nil {
		r.Error = c.DefaultError
	}
	//
	// Create a URL object from the raw url string.  This will allow us to compose
	// query parameters programmatically and be guaranteed of a well-formed URL.
	//
	u, err := url.Parse(r.Url)
	if err != nil {
		log.Println(err)
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
		if err != nil {
			log.Println(err)
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(m, u.String(), buf)
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		log.Println(err)
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
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		lineNo := strconv.Itoa(line)
		s := "Error executing REST request, called from " + file + ":" + lineNo + ": "
		log.Println(s, err)
		return
	}
	status = resp.StatusCode
	r.Status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	r.RawText = string(data)
	// If server returned no data, don't bother trying to unmarshall it (which will fail anyways).
	if r.RawText == "" {
		return
	}
	if status >= 200 && status < 300 {
		err = c.unmarshal(data, &r.Result)
	} else {
		err = c.unmarshal(data, &r.Error)
	}
	if err != nil {
		log.Println(status)
		log.Println(err)
		log.Println(r.RawText)
	}
	return
}

// unmarshal parses the JSON-encoded data and stores the result in the value
// pointed to by v.  If the data cannot be unmarshalled without error, v will be 
// reassigned the value interface{}, and data unmarshalled into that.
func (c *Client) unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err == nil {
		return nil
	}
	v = new(interface{})
	return json.Unmarshal(data, v)
}
