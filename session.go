// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

/*
This module provides a Session object to manage and persist settings across
requests (cookies, auth, proxies).
*/

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

import ()

type Session struct {
	Opts            *Opts
	Client          *http.Client
	UnsafeBasicAuth bool // Allow Basic Auth over unencrypted HTTP
	Log             bool // Log request and response
}

// Send constructs and sends an HTTP request.
func (s *Session) Send(r *Request) (status int, err error) {
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
		for k, v := range *r.Params {
			vals.Set(k, v)
		}
		u.RawQuery = vals.Encode()
	}
	//
	// Create a Request object; if populated, Data field is JSON encoded as
	// request body
	//
	var req *http.Request
	if r.Data != nil {
		var b []byte
		b, err = json.Marshal(&r.Data)
		if err != nil {
			log.Println(err)
			return
		}
		buf := bytes.NewBuffer(b)
		req, err = http.NewRequest(r.Method, u.String(), buf)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
	} else { // no data to encode
		req, err = http.NewRequest(r.Method, u.String(), nil)
		if err != nil {
			log.Println(err)
			return
		}

	}
	req.Header.Add("Accept", "application/json") // Default, can be overridden with Opts
	o := s.Opts.update(r.Opts)
	if o.Header != nil {
		for key, values := range *o.Header {
			if len(values) > 0 {
				req.Header.Set(key, values[0]) // Possible to overwrite Content-Type
			}
		}
	}
	//
	// Set HTTP Basic authentication if userinfo is supplied
	//
	if o.Userinfo != nil {
		if !s.UnsafeBasicAuth && u.Scheme != "https" {
			err = errors.New("Unsafe to use HTTP Basic authentication without HTTPS")
			return
		}
		pwd, _ := o.Userinfo.Password()
		req.SetBasicAuth(o.Userinfo.Username(), pwd)
	}
	//
	// Execute the HTTP request
	//
	if s.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("REQUEST")
		log.Println("--------------------------------------------------------------------------------")
		prettyPrint(req)
		log.Print("Payload: ")
		prettyPrint(r.Data)
	}
	r.timestamp = time.Now()
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	r.status = resp.StatusCode
	r.response = resp
	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(r.body, &r.Result)
	if s.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("RESPONSE")
		log.Println("--------------------------------------------------------------------------------")
		log.Println("Status: ", r.status)
		if r.body != nil {
			raw := json.RawMessage{}
			if json.Unmarshal(r.body, &raw) == nil {
				prettyPrint(&raw)
			} else {
				prettyPrint(r.RawText)
			}
		} else {
			log.Println("Empty response body")
		}

	}
	if o.ExpectedStatus != 0 && r.status != o.ExpectedStatus {
		log.Printf("Expected status %s but got %s", o.ExpectedStatus, r.status)
		return status, UnexpectedStatus
	}
	return
}

// Get sends a GET request.
func (s *Session) Get(url string, p *Params, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Options sends an OPTIONS request.
func (s *Session) Options(url string, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Head sends a HEAD request.
func (s *Session) Head(url string, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Post sends a POST request.
func (s *Session) Post(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Put sends a PUT request.
func (s *Session) Put(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Patch sends a PATCH request.
func (s *Session) Patch(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Delete sends a DELETE request.
func (s *Session) Delete(url string, o *Opts) (status int, err error) {
	return 0, nil
}
