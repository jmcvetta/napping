// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

/*
This module implements the Napping API.
*/

import ()

// Get sends a GET request.
func Get(url string, p *Params, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Options sends an OPTIONS request.
func Options(url string, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Head sends a HEAD request.
func Head(url string, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Post sends a POST request.
func Post(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Put sends a PUT request.
func Put(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Patch sends a PATCH request.
func Patch(url string, data, response interface{}, o *Opts) (status int, err error) {
	return 0, nil
}

// Delete sends a DELETE request.
func Delete(url string, o *Opts) (status int, err error) {
	return 0, nil
}
