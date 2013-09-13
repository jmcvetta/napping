// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A Params is a map containing URL parameters.
type Params map[string]string

// A UnexpectedStatus error is returned when ExpectedStatus is set, and the
// server return a status code other than what is expected.
var UnexpectedStatus = errors.New("Server returned unexpected status.")

// A RequestResponse describes an HTTP request to be executed, data
// structures into which results and errors will be unmarshalled, and the
// server's response.  By using a single object for both the request and the
// response we allow easy access to Result and Error objects without needing
// type assertions.
type RequestResponse struct {
	Url      string        // Raw URL string
	Method   string        // HTTP method to use
	Userinfo *url.Userinfo // Optional username/password to authenticate this request
	Params   Params        // URL parameters for GET requests (ignored otherwise)
	Header   *http.Header  // HTTP Headers to use (will override defaults)
	//
	// OPTIONAL - ExpectedStatus is the HTTP status code we expect the server
	// to return on a successful request.  If ExpectedStatus is non-zero and
	// server returns a different code, Client.Do will return a BadStatus error.
	//
	ExpectedStatus int
	//
	// The following interfaces fields should be populated with *pointers* to
	// data structures.  Any structure that can be (un)marshalled by the json
	// package can be used.
	//
	Data   interface{} // Data to JSON-encode and POST
	Result interface{} // Successful response is unmarshalled into Result
	Error  interface{} // Error response is unmarshalled into Error
	//
	// The following fields are populated by Client.Do()
	//
	Timestamp    time.Time      // Time when HTTP request was sent
	RawText      string         // Raw text of server response (JSON or otherwise)
	Status       int            // HTTP status for executed request
	HttpResponse *http.Response // Response object from http package
}

// Client is a REST client.
type Client struct {
	HttpClient      *http.Client
	UnsafeBasicAuth bool // Allow Basic Auth over unencrypted HTTP
	Log             bool // Log request and response
}

// New returns a new Client instance.
func New() *Client {
	return &Client{
		HttpClient:      new(http.Client),
		UnsafeBasicAuth: false,
	}
}

// Do executes a REST request.
func (c *Client) Do(rr *RequestResponse) (status int, err error) {
	rr.Method = strings.ToUpper(rr.Method)
	//
	// Create a URL object from the raw url string.  This will allow us to compose
	// query parameters programmatically and be guaranteed of a well-formed URL.
	//
	u, err := url.Parse(rr.Url)
	if err != nil {
		log.Println(err)
		return
	}
	//
	// If we are making a GET request and the user populated the Params field, then
	// add the params to the URL's querystring.
	//
	if rr.Method == "GET" && rr.Params != nil {
		vals := u.Query()
		for k, v := range rr.Params {
			vals.Set(k, v)
		}
		u.RawQuery = vals.Encode()
	}
	//
	// Create a Request object; if populated, Data field is JSON encoded as request
	// body
	//
	rr.Timestamp = time.Now()
	m := string(rr.Method)
	var req *http.Request
	// http.NewRequest can only return an error if url.Parse fails.  Since the
	// url has already been successfully parsed once at this point, there is no
	// danger of this, so we can ignore errors returned by http.NewRequest.
	if rr.Data == nil {
		req, _ = http.NewRequest(m, u.String(), nil)
	} else {
		var b []byte
		b, err = json.Marshal(&rr.Data)
		if err != nil {
			log.Println(err)
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(m, u.String(), buf)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
	}
	if rr.Header != nil {
		for key, values := range *rr.Header {
			if len(values) > 0 {
				req.Header.Set(key, values[0]) // Possible to overwrite Content-Type
			}
		}
	}
	//
	// If Accept header is unset, set it for JSON.
	//
	if req.Header.Get("Accept") == "" {
		req.Header.Add("Accept", "application/json")
	}
	//
	// Set HTTP Basic authentication if userinfo is supplied
	//
	if rr.Userinfo != nil {
		if !c.UnsafeBasicAuth && u.Scheme != "https" {
			err = errors.New("Unsafe to use HTTP Basic authentication without HTTPS")
			return
		}
		pwd, _ := rr.Userinfo.Password()
		req.SetBasicAuth(rr.Userinfo.Username(), pwd)
	}
	//
	// Execute the HTTP request
	//
	if c.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("REQUEST")
		log.Println("--------------------------------------------------------------------------------")
		prettyPrint(req)
		log.Print("Payload: ")
		prettyPrint(rr.Data)
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	rr.HttpResponse = resp
	rr.Status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	rr.RawText = string(data)
	json.Unmarshal(data, &rr.Error) // Ignore errors
	if rr.RawText != "" && status < 300 {
		err = json.Unmarshal(data, &rr.Result) // Ignore errors
	}
	if c.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("RESPONSE")
		log.Println("--------------------------------------------------------------------------------")
		log.Println("Status: ", status)
		if rr.RawText != "" {
			raw := json.RawMessage{}
			if json.Unmarshal(data, &raw) == nil {
				prettyPrint(&raw)
			} else {
				prettyPrint(rr.RawText)
			}
		} else {
			log.Println("Empty response body")
		}

	}
	if rr.ExpectedStatus != 0 && status != rr.ExpectedStatus {
		log.Printf("Expected status %s but got %s", rr.ExpectedStatus, status)
		return status, UnexpectedStatus
	}
	return
}

var (
	defaultClient = New()
)

// Do executes a REST request using the default client.
func Do(rr *RequestResponse) (status int, err error) {
	return defaultClient.Do(rr)
}
