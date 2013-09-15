// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

/*
This module provides a Session object to manage and persist settings across
requests (cookies, auth, proxies).
*/

import ()

type Session struct {
}

// Request constructs and sends an HTTP request.
func (s *Session) Send(r *Request) (status int, err error) {
	return 0, nil
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
