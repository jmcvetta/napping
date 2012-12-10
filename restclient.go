// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package restclient provides a simple client library for interacting with
// RESTful APIs.
package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RestClient struct {
	client *http.Client
}

func (r *RestClient) GetUrl(u *url.URL, result interface{}) (status int, err error) {
	c := call{
		Url:    u,
		Method: "GET",
		Result: result,
	}
	return r.rest(&c)
}

func (r *RestClient) Get(rawurl string, result interface{}) (status int, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}
	return r.GetUrl(u, result)
}

type call struct {
	Url     *url.URL    // Absolute URL to call
	Method  string      // HTTP method to use 
	Content interface{} // Data to JSON-encode and include with call
	Result  interface{} // JSON-encoded data in respose will be unmarshalled into Result
}

func (r *RestClient) rest(c *call) (status int, err error) {
	req, err := http.NewRequest(c.Method, c.Url.String(), nil)
	if err != nil {
		return
	}
	if c.Content != nil {
		// log.Println(pretty.Sprintf("Content: %# v", c.Content))
		var b []byte
		b, err = json.Marshal(c.Content)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(c.Method, c.Url.String(), buf)
		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Accept", "application/json")
	// log.Println(pretty.Sprintf("Request: %# v", req))
	resp, err := r.client.Do(req)
	if err != nil {
		return
	}
	status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	// Ignore unmarshall errors - worst case is, c.Result will be nil
	err = json.Unmarshal(data, &c.Result)
	return
}
