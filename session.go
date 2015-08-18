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
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Session struct {
	EncodingMarshaller
	Client *http.Client
	Log    bool // Log request and response

	// Optional
	Userinfo *url.Userinfo

	// Optional defaults - can be overridden in a Request
	Header *http.Header
	Params *Params
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
		s.log("URL", r.Url)
		s.log(err)
		return
	}
	//
	// Default query parameters
	//
	p := Params{}
	if s.Params != nil {
		for k, v := range *s.Params {
			p[k] = v
		}
	}
	//
	// User-supplied params override default
	//
	if r.Params != nil {
		for k, v := range *r.Params {
			p[k] = v
		}
	}
	//
	// Encode parameters
	//
	vals := u.Query()
	for k, v := range p {
		vals.Set(k, v)
	}
	u.RawQuery = vals.Encode()
	//
	// Create a Request object; if populated, Data field is encoded as
	// request body
	//
	header := http.Header{}
	if s.Header != nil {
		for k, _ := range *s.Header {
			v := s.Header.Get(k)
			header.Set(k, v)
		}
	}
	var req *http.Request
	var buf *bytes.Buffer
	if r.Payload != nil {
		if r.RawPayload {
			var ok bool
			// buf can be nil interface at this point
			// so we'll do extra nil check
			buf, ok = r.Payload.(*bytes.Buffer)
			if !ok {
				err = errors.New("Payload must be of type *bytes.Buffer if RawPayload is set to true")
				return
			}
		} else {
			var b []byte
			b, err = s.Marshal(&r.Payload)
			if err != nil {
				s.log(err)
				return
			}
			buf = bytes.NewBuffer(b)
		}
		if buf != nil {
			req, err = http.NewRequest(r.Method, u.String(), buf)
		} else {
			req, err = http.NewRequest(r.Method, u.String(), nil)
		}
		if err != nil {
			s.log(err)
			return
		}
		header.Add("Content-Type", s.getContentType())
	} else { // no data to encode
		req, err = http.NewRequest(r.Method, u.String(), nil)
		if err != nil {
			s.log(err)
			return
		}

	}
	//
	// Merge Session and Request options
	//
	var userinfo *url.Userinfo
	if u.User != nil {
		userinfo = u.User
	}
	if s.Userinfo != nil {
		userinfo = s.Userinfo
	}
	// Prefer Request's user credentials
	if r.Userinfo != nil {
		userinfo = r.Userinfo
	}
	if r.Header != nil {
		for k, v := range *r.Header {
			header.Set(k, v[0]) // Is there always guarnateed to be at least one value for a header?
		}
	}
	if header.Get("Accept") == "" {
		header.Add("Accept", s.getContentType()) // Default, can be overridden with Opts
	}
	req.Header = header
	//
	// Set HTTP Basic authentication if userinfo is supplied
	//
	if userinfo != nil {
		pwd, _ := userinfo.Password()
		req.SetBasicAuth(userinfo.Username(), pwd)
		if u.Scheme != "https" {
			s.log("WARNING: Using HTTP Basic Auth in cleartext is insecure.")
		}
	}
	//
	// Execute the HTTP request
	//

	// Debug log request
	s.log("--------------------------------------------------------------------------------")
	s.log("REQUEST")
	s.log("--------------------------------------------------------------------------------")
	s.logWithContext(req)
	s.log("Payload:")
	if r.RawPayload && s.Log && buf != nil {
		s.log(base64.StdEncoding.EncodeToString(buf.Bytes()))
	} else {
		s.logWithContext(r.Payload)
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
		s.log(err)
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
		s.log(err)
		return
	}
	if string(r.body) != "" {
		if resp.StatusCode < 300 && r.Result != nil {
			err = s.Unmarshal(r.body, r.Result)
		}
		if resp.StatusCode >= 400 && r.Error != nil {
			_ = s.Unmarshal(r.body, r.Error) // Should we ignore unmarshall error?
		}
	}
	rsp := Response(*r)
	response = &rsp

	// Debug log response
	s.log("--------------------------------------------------------------------------------")
	s.log("RESPONSE")
	s.log("--------------------------------------------------------------------------------")
	s.log("Status: ", response.status)
	s.log("Header:")
	s.logWithContext(response.HttpResponse().Header)
	s.log("Body:")

	if response.body != nil {
		s.logRaw(response)
	} else {
		s.log("Empty response body")
	}

	return
}

func (s *Session) getContentType() string {
	switch s.Encoding {
	case JSON:
		return "application/json"
	case XML:
		return "application/xml"
	}
	panic("invalid encoding")
}

// Get sends a GET request.
func (s *Session) Get(url string, p *Params, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:             "GET",
		Url:                url,
		Params:             p,
		Result:             result,
		Error:              errMsg,
		EncodingMarshaller: EncodingMarshaller{Encoding: s.Encoding},
	}
	return s.Send(&r)
}

// Options sends an OPTIONS request.
func (s *Session) Options(url string, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method: "OPTIONS",
		Url:    url,
		Result: result,
		Error:  errMsg,
	}
	return s.Send(&r)
}

// Head sends a HEAD request.
func (s *Session) Head(url string, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method: "HEAD",
		Url:    url,
		Result: result,
		Error:  errMsg,
	}
	return s.Send(&r)
}

// Post sends a POST request.
func (s *Session) Post(url string, payload, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:  "POST",
		Url:     url,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}
	return s.Send(&r)
}

// Put sends a PUT request.
func (s *Session) Put(url string, payload, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:  "PUT",
		Url:     url,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}
	return s.Send(&r)
}

// Patch sends a PATCH request.
func (s *Session) Patch(url string, payload, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method:  "PATCH",
		Url:     url,
		Payload: payload,
		Result:  result,
		Error:   errMsg,
	}
	return s.Send(&r)
}

// Delete sends a DELETE request.
func (s *Session) Delete(url string, result, errMsg interface{}) (*Response, error) {
	r := Request{
		Method: "DELETE",
		Url:    url,
		Result: result,
		Error:  errMsg,
	}
	return s.Send(&r)
}

// Debug method for logging
// Centralizing logging in one method
// avoids spreading conditionals everywhere
func (s *Session) log(args ...interface{}) {
	if s.Log {
		log.Println(args...)
	}
}

// log the raw text of a response
func (s *Session) logRaw(response *Response) {
	switch s.Encoding {
	case JSON:
		s.logWithContext(response.RawText())
		return
	case XML:
		s.logWithContext(response.RawText())
		return
	}
	panic("invalid encoding")
}

// log a Go representation of a value, with caller context information
func (s *Session) logWithContext(v interface{}) {
	if s.Log {
		// Get source file and line
		_, file, line, _ := runtime.Caller(1)

		// Get relative filename
		filename := filepath.Base(file)

		// Make that all pretty together
		s.log(fmt.Sprintf("%s:%d: \n%+v\n", filename, line, v))
	}
}
