// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

/*
Package restclient is a client library for interacting with RESTful APIs.

Example:

	type Foo struct {
		Bar string
	}
	type Spam struct {
		Eggs int
	}
	f := Foo{
		Bar: "baz",
	}
	s := Spam{}
	r := restclient.RequestResponse{
		Url:    "http://foo.com/bar",
		Method: restclient.POST,
		Data:   &f,
		Result: &s,
	}
	status, err := restclient.Do(&r)
	if err != nil {
		panic(err)
	}
	if status == 200 {
		println(s.Eggs)
	}
*/
package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// A Method is an HTTP verb.
type Method string

type Params map[string]string

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
	Timestamp time.Time // Time when HTTP request was sent
	RawText   string    // Raw text of server response (JSON or otherwise)
	Status    int       // HTTP status for executed request
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
func (c *Client) Do(r *RequestResponse) (status int, err error) {
	r.Method = strings.ToUpper(r.Method)
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
	if r.Method == "GET" && r.Params != nil {
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
	r.Timestamp = time.Now()
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
	if r.Header != nil {
		for key, values := range *r.Header {
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
	if r.Userinfo != nil {
		if !c.UnsafeBasicAuth && u.Scheme != "https" {
			err = errors.New("Unsafe to use HTTP Basic authentication without HTTPS")
			return
		}
		pwd, _ := r.Userinfo.Password()
		req.SetBasicAuth(r.Userinfo.Username(), pwd)
	}
	//
	// Execute the HTTP request
	//
	if c.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("REQUEST")
		log.Println("--------------------------------------------------------------------------------")
		b, _ := json.MarshalIndent(req, "", "\t")
		log.Println("\n" + string(b))
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		complain(err, status, "")
		return
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	r.Status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		complain(err, status, string(data))
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
		log.Println(resp)
		log.Println(resp.Request)
	}
	if c.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("RESPONSE")
		log.Println("--------------------------------------------------------------------------------")
		log.Println("Status: ", status)
		raw := json.RawMessage{}
		c.unmarshal(data, &raw)
		b, _ := json.MarshalIndent(&raw, "", "\t")
		log.Println("\n" + string(b))
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

// complain prints detailed error messages to the log.
func complain(err error, status int, rawtext string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	lineNo := strconv.Itoa(line)
	s := "Error executing REST request:\n"
	s += "    --> Called from " + file + ":" + lineNo + "\n"
	s += "    --> Got status " + strconv.Itoa(status) + "\n"
	if rawtext != "" {
		s += "    --> Raw text of server response: " + rawtext + "\n"
	}
	s += "    --> " + err.Error()
	log.Println(s)
}

var (
	defaultClient = New()
)

// Do executes a REST request using the default client.
func Do(r *RequestResponse) (status int, err error) {
	return defaultClient.Do(r)
}
