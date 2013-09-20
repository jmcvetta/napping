// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import "io/ioutil"

import "errors"

import "bytes"

import "encoding/json"

import "net/http"

import "time"

import "net/url"

import "log"

import "strings"

/*
This module provides a Session object to manage and persist settings across
requests (cookies, auth, proxies).
*/

import ()

type Session struct {
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
	// Create a Request object; if populated, Data field is JSON encoded as request
	// body
	//
	r.Timestamp = time.Now()
	m := string(r.Method)
	var req *http.Request
	// http.NewRequest can only return an error if url.Parse fails.  Since the
	// url has already been successfully parsed once at this point, there is no
	// danger of this, so we can ignore errors returned by http.NewRequest.
	if r.Data == nil {
		req, _ = http.NewRequest(m, u.String(), nil)
	} else {
		var b []byte
		b, err = json.Marshal(&r.Data)
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
		if !s.UnsafeBasicAuth && u.Scheme != "https" {
			err = errors.New("Unsafe to use HTTP Basic authentication without HTTPS")
			return
		}
		pwd, _ := r.Userinfo.Password()
		req.SetBasicAuth(r.Userinfo.Username(), pwd)
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
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	r.HttpResponse = resp
	r.Status = resp.StatusCode
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	r.RawText = string(data)
	json.Unmarshal(data, &r.Error) // Ignore errors
	if r.RawText != "" && status < 300 {
		err = json.Unmarshal(data, &r.Result) // Ignore errors
	}
	if s.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("RESPONSE")
		log.Println("--------------------------------------------------------------------------------")
		log.Println("Status: ", status)
		if r.RawText != "" {
			raw := json.RawMessage{}
			if json.Unmarshal(data, &raw) == nil {
				prettyPrint(&raw)
			} else {
				prettyPrint(r.RawText)
			}
		} else {
			log.Println("Empty response body")
		}

	}
	if r.ExpectedStatus != 0 && status != r.ExpectedStatus {
		log.Printf("Expected status %s but got %s", r.ExpectedStatus, status)
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
