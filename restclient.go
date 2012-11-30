// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

// Package restclient provides a simple client library for interacting with
// RESTful APIs.
package restclient

import (
	"net/url"
	"net/http"
	"encoding/json"
	"bytes"
)

type RestClient struct {
	client *http.Client
}

func (rc *RestClient) Get(u url.URL, Result interface{}) (status int, err error) {
}

type restCall struct {
	Url     url.URL      // Absolute URL to call
	Method  string      // HTTP method to use 
	Content interface{} // Data to JSON-encode and include with call
	Result  interface{} // JSON-encoded data in respose will be unmarshalled into Result
}

func rest(r *restCall) (status int, err error) {
	req, err := http.NewRequest(r.Method, r.Url.String(), nil)
	if err != nil {
		return
	}
	if r.Content != nil {
		// log.Println(pretty.Sprintf("Content: %# v", r.Content))
		var b []byte
		b, err = json.Marshal(r.Content)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(r.Method, r.Url.String(), buf)
		if err != nil {
			return
		}
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Accept", "application/json")
	// log.Println(pretty.Sprintf("Request: %# v", req))
	resp, err := db.client.Do(req)
	if err != nil {
		return
	}
	status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	// Ignore unmarshall errors - worst case is, r.Result will be nil
	json.Unmarshal(data, &r.Result)
	if status < 200 || status >= 300 {
		res := &r.Result
		// log.Println(*res)
		info, ok := (*res).(neoError)
		if ok {
			log.Println("Got error response code:", status)
			log.Println(info.Mesage)
			log.Println(info.Exception)
			log.Println(info.StackTrace)
		}
	}
	return
}
