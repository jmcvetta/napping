// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.M06.  http://docs.neo4j.org/chunked/milestone/

package restclient

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type structType struct {
	Foo int
	Bar string
}

type errorStruct struct {
	Status  int
	Message string
}

var (
	fooMap    = map[string]string{"foo": "bar"}
	barMap    = map[string]string{"bar": "baz"}
	fooStruct = structType{
		Foo: 111,
		Bar: "foo",
	}
	barStruct = structType{
		Foo: 222,
		Bar: "bar",
	}
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func JsonError(w http.ResponseWriter, msg string, code int) {
	e := errorStruct{
		Status:  code,
		Message: msg,
	}
	blob, err := json.Marshal(e)
	if err != nil {
		http.Error(w, msg, code)
		return
	}
	http.Error(w, string(blob), code)
}

func HandleGet(w http.ResponseWriter, req *http.Request) {
	u := req.URL
	q := u.Query()
	for k, _ := range fooMap {
		if fooMap[k] != q.Get(k) {
			msg := "Bad query params: " + u.Query().Encode()
			JsonError(w, msg, http.StatusInternalServerError)
			return
		}
	}
	//
	// Generate response
	//
	blob, err := json.Marshal(barStruct)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(blob)
}
func HandlePost(w http.ResponseWriter, req *http.Request) {
	//
	// Parse Payload
	//
	if req.ContentLength <= 0 {
		msg := "Content-Length must be greater than 0."
		JsonError(w, msg, http.StatusLengthRequired)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var s structType
	err = json.Unmarshal(body, &s)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s != fooStruct {
		msg := "Bad request body"
		JsonError(w, msg, http.StatusBadRequest)
		return
	}
	//
	// Compose Response
	//
	blob, err := json.Marshal(barStruct)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(blob)
}

func HandlePut(w http.ResponseWriter, req *http.Request) {
	//
	// Parse Payload
	//
	if req.ContentLength <= 0 {
		msg := "Content-Length must be greater than 0."
		JsonError(w, msg, http.StatusLengthRequired)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var s structType
	err = json.Unmarshal(body, &s)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s != fooStruct {
		msg := "Bad request body"
		JsonError(w, msg, http.StatusBadRequest)
		return
	}
	return
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGet))
	defer srv.Close()
	// 
	// Good request
	//
	client := New()
	r := RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: GET,
		Params: fooMap,
		// Params: map[string]string{"bad": "value"},
		Result: new(structType),
	}
	status, err := client.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	assert.Equal(t, r.Result, &barStruct)
	// 
	// Bad request
	//
	client = New()
	r = RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: GET,
		Params: map[string]string{"bad": "value"},
		Error:  new(errorStruct),
	}
	status, err = client.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 500)
	expected := errorStruct{
		Message: "Bad query params: bad=value",
		Status:  500,
	}
	e := r.Error.(*errorStruct)
	assert.Equal(t, *e, expected)
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePost))
	defer srv.Close()
	client := New()
	r := RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: POST,
		Data:   fooStruct,
		Result: new(structType),
	}
	status, err := client.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	assert.Equal(t, r.Result, &barStruct)
}

func TestPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePut))
	defer srv.Close()
	client := New()
	r := RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: PUT,
		Data:   fooStruct,
		Result: new(structType),
	}
	status, err := client.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	// Server should return NO data
	assert.Equal(t, r.RawText, "")
}
