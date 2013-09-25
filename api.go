// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

/*
This module implements the Napping API.
*/

import ()

// Send composes and sends and HTTP request.
func Send(r *Request) (response *Response, err error) {
	s := Session{}
	return s.Send(r)
}

// Get sends a GET request.
func Get(url string, p *Params, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Get(url, p, result, o)
}

// Options sends an OPTIONS request.
func Options(url string, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Options(url, result, o)
}

// Head sends a HEAD request.
func Head(url string, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Head(url, result, o)
}

// Post sends a POST request.
func Post(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Post(url, payload, result, o)
}

// Put sends a PUT request.
func Put(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Put(url, payload, result, o)
}

// Patch sends a PATCH request.
func Patch(url string, payload, result interface{}, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Patch(url, payload, result, o)
}

// Delete sends a DELETE request.
func Delete(url string, o *Opts) (response *Response, err error) {
	s := Session{}
	return s.Delete(url, o)
}
