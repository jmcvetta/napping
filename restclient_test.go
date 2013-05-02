// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package restclient

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/randutil"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

//
// Request Tests
//

type hfunc http.HandlerFunc

type payload struct {
	Foo string
}

var reqTests = []struct {
	method  string
	params  bool
	payload bool
}{
	{"GET", true, false},
	{"POST", false, true},
	{"PUT", false, true},
	{"DELETE", false, false},
}

type pair struct {
	rr RequestResponse
	hf hfunc
}

func paramHandler(t *testing.T, p Params, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		q := req.URL.Query()
		for k, _ := range p {
			if p[k] != q.Get(k) {
				msg := "Bad query params: " + q.Encode()
				t.Error(msg)
				return
			}
		}
	}
}

func payloadHandler(t *testing.T, p payload, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		if req.ContentLength <= 0 {
			t.Error("Content-Length must be greater than 0.")
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("Bad content type")
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Error("Body is nil")
			return
		}
		var s payload
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Error("JSON Unmarshal failed: ", err)
			return
		}
		if s != p {
			t.Error("Bad request body")
			return
		}
	}
}

func methodHandler(t *testing.T, method string, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		if req.Method != method {
			t.Error("Incorrect method, got ", req.Method, " expected ", method)
		}
	}
}

func headerHandler(t *testing.T, h http.Header, f hfunc) hfunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if f != nil {
			f(w, req)
		}
		for key, _ := range h {
			expected := h.Get(key)
			actual := req.Header.Get(key)
			if expected != actual {
				t.Error("Missing/bad header")
			}
			return
		}
	}
}

func TestRequest(t *testing.T) {
	// NOTE:  Do we really need to test different combinations for different
	// HTTP methods?
	pairs := []pair{}
	for _, test := range reqTests {
		baseRR := RequestResponse{
			Method: test.method,
		}
		allRR := baseRR // allRR has all supported attribues for this verb
		var allHF hfunc // allHF is combination of all relevant handlers
		//
		// Generate a random key/value pair
		//
		key, err := randutil.AlphaString(8)
		if err != nil {
			t.Error(err)
		}
		value, err := randutil.AlphaString(8)
		if err != nil {
			t.Error(err)
		}
		//
		// Method
		//
		r := baseRR
		f := methodHandler(t, test.method, nil)
		allHF = methodHandler(t, test.method, allHF)
		pairs = append(pairs, pair{r, f})
		//
		// Header
		//
		h := http.Header{}
		h.Add(key, value)
		r = baseRR
		r.Header = &h
		allRR.Header = &h
		f = headerHandler(t, h, nil)
		allHF = headerHandler(t, h, allHF)
		pairs = append(pairs, pair{r, f})
		//
		// Params
		//
		if test.params {
			p := Params{key: value}
			f := paramHandler(t, p, nil)
			allHF = paramHandler(t, p, allHF)
			r = baseRR
			r.Params = p
			allRR.Params = p
			pairs = append(pairs, pair{r, f})
		}
		//
		// Payload
		//
		if test.payload {
			p := payload{value}
			f = payloadHandler(t, p, nil)
			allHF = payloadHandler(t, p, allHF)
			r = baseRR
			r.Data = p
			allRR.Data = p
			pairs = append(pairs, pair{r, f})
		}
		//
		// All
		//
		pairs = append(pairs, pair{allRR, allHF})
	}
	for _, p := range pairs {
		srv := httptest.NewServer(http.HandlerFunc(p.hf))
		defer srv.Close()
		//
		// Good request
		//
		client := New()
		p.rr.Url = "http://" + srv.Listener.Addr().String()
		_, err := client.Do(&p.rr)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidUrl(t *testing.T) {
	client := New()
	rr := RequestResponse{
		Url:    "carrots",
		Method: "GET",
	}
	_, err := client.Do(&rr)
	assert.NotEqual(t, nil, err)

}

//
// TODO: Response Tests
//

func TestErrMsg(t *testing.T) {}

func TestStatus(t *testing.T) {}

func TestUnmarshall(t *testing.T) {}

// func TestUnmarshallFail() {}

//
// Old Tests
//

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

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGet))
	defer srv.Close()
	//
	// Good request
	//
	client := New()
	r := RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: "GET",
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
		Method: "GET",
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
	client.Log = true
	r := RequestResponse{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: "POST",
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
		Method: "PUT",
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
