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
func (s *Session) Send(r *Request) (response *Response, err error) {
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
	if r.Payload != nil {
		var b []byte
		b, err = json.Marshal(&r.Payload)
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
	sessOpts := s.Opts
	if sessOpts == nil {
		sessOpts = &Opts{}
	}
	mergedOpts := sessOpts.merge(r.Opts)
	if mergedOpts.Header != nil {
		for key, values := range *mergedOpts.Header {
			if len(values) > 0 {
				req.Header.Set(key, values[0]) // Possible to overwrite Content-Type
			}
		}
	}
	//
	// Set HTTP Basic authentication if userinfo is supplied
	//
	if mergedOpts.Userinfo != nil {
		if !s.UnsafeBasicAuth && u.Scheme != "https" {
			err = errors.New("Unsafe to use HTTP Basic authentication without HTTPS")
			return
		}
		pwd, _ := mergedOpts.Userinfo.Password()
		req.SetBasicAuth(mergedOpts.Userinfo.Username(), pwd)
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
		prettyPrint(r.Payload)
	}
	r.timestamp = time.Now()
	var client *http.Client
	if s.Client != nil {
		client = s.Client
	} else {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	r.status = resp.StatusCode
	r.response = resp
	//
	// Unmarshal
	//
	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode < 300 && r.Result != nil && string(r.body) != "" {
		err = json.Unmarshal(r.body, r.Result)
	}
	rsp := Response(*r)
	response = &rsp
	if s.Log {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("RESPONSE")
		log.Println("--------------------------------------------------------------------------------")
		log.Println("Status: ", response.status)
		if response.body != nil {
			raw := json.RawMessage{}
			if json.Unmarshal(response.body, &raw) == nil {
				prettyPrint(&raw)
			} else {
				prettyPrint(response.RawText())
			}
		} else {
			log.Println("Empty response body")
		}

	}
	if mergedOpts.ExpectedStatus != 0 && r.status != mergedOpts.ExpectedStatus {
		log.Printf("Expected status %v but got %v", mergedOpts.ExpectedStatus, r.status)
		return response, UnexpectedStatus
	}
	return
}

// Get sends a GET request.
func (s *Session) Get(url string, p *Params, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method: "GET",
		Url:    url,
		Params: p,
		Opts:   o,
		Result: result,
	}
	return s.Send(&r)
}

// Options sends an OPTIONS request.
func (s *Session) Options(url string, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method: "OPTIONS",
		Url:    url,
		Opts:   o,
		Result: result,
	}
	return s.Send(&r)
}

// Head sends a HEAD request.
func (s *Session) Head(url string, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method: "HEAD",
		Url:    url,
		Opts:   o,
		Result: result,
	}
	return s.Send(&r)
}

// Post sends a POST request.
func (s *Session) Post(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method:  "POST",
		Url:     url,
		Opts:    o,
		Payload: payload,
		Result:  result,
	}
	return s.Send(&r)
}

// Put sends a PUT request.
func (s *Session) Put(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method:  "PUT",
		Url:     url,
		Opts:    o,
		Payload: payload,
		Result:  result,
	}
	return s.Send(&r)
}

// Patch sends a PATCH request.
func (s *Session) Patch(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	r := Request{
		Method:  "PATCH",
		Url:     url,
		Opts:    o,
		Payload: payload,
		Result:  result,
	}
	return s.Send(&r)
}

// Delete sends a DELETE request.
func (s *Session) Delete(url string, o *Opts) (response *Response, err error) {
	r := Request{
		Method: "DELETE",
		Url:    url,
		Opts:   o,
	}
	return s.Send(&r)
}
