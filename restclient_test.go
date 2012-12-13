// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

//
// The Neo4j Manual section numbers quoted herein refer to the manual for 
// milestone release 1.8.M06.  http://docs.neo4j.org/chunked/milestone/

package restclient

import (
	"github.com/bmizerany/assert"
	"github.com/bmizerany/pat"
	"log"
	// "sort"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"testing"
)

const (
	port = "9000"
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
	// 
	// Routing
	//
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(HandleGET))
	mux.Post("/", http.HandlerFunc(HandlePOST))
	//
	// Start webserver
	//
	http.Handle("/", mux)
	go func() {
		log.Println("Starting webserver on port " + port + "...")
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Panicln(err)
		}
	}()
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

func HandleGET(w http.ResponseWriter, req *http.Request) {
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

func HandlePOST(w http.ResponseWriter, req *http.Request) {
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

func TestGET(t *testing.T) {
	// 
	// Good request
	//
	c := New()
	r := RestRequest{
		Url:    "http://localhost:" + port,
		Method: GET,
		Params: fooMap,
		// Params: map[string]string{"bad": "value"},
		Result: new(structType),
	}
	status, err := c.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	assert.Equal(t, r.Result, &barStruct)
	// 
	// Bad request
	//
	c = New()
	r = RestRequest{
		Url:    "http://localhost:" + port,
		Method: GET,
		Params: map[string]string{"bad": "value"},
		Error:  new(errorStruct),
	}
	status, err = c.Do(&r)
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

func TestPOST(t *testing.T) {
	c := New()
	r := RestRequest{
		Url:    "http://localhost:" + port,
		Method: POST,
		Data:   fooStruct,
		Result: new(structType),
	}
	status, err := c.Do(&r)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, status, 200)
	assert.Equal(t, r.Result, &barStruct)
}
