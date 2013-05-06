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
