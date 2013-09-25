// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"encoding/base64"
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/randutil"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
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
	r  Request
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
		baseReq := Request{
			Method: test.method,
			Opts:   &Opts{},
		}
		allReq := baseReq // allRR has all supported attribues for this verb
		var allHF hfunc   // allHF is combination of all relevant handlers
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
		r := baseReq
		f := methodHandler(t, test.method, nil)
		allHF = methodHandler(t, test.method, allHF)
		pairs = append(pairs, pair{r, f})
		//
		// Header
		//
		h := http.Header{}
		h.Add(key, value)
		r = baseReq
		r.Opts.Header = &h
		allReq.Opts.Header = &h
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
			r = baseReq
			r.Params = &p
			allReq.Params = &p
			pairs = append(pairs, pair{r, f})
		}
		//
		// Payload
		//
		if test.payload {
			p := payload{value}
			f = payloadHandler(t, p, nil)
			allHF = payloadHandler(t, p, allHF)
			r = baseReq
			r.Payload = p
			allReq.Payload = p
			pairs = append(pairs, pair{r, f})
		}
		//
		// All
		//
		pairs = append(pairs, pair{allReq, allHF})
	}
	for _, p := range pairs {
		srv := httptest.NewServer(http.HandlerFunc(p.hf))
		defer srv.Close()
		//
		// Good request
		//
		p.r.Url = "http://" + srv.Listener.Addr().String()
		_, err := Send(&p.r)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidUrl(t *testing.T) {
	//
	//  Missing protocol scheme - url.Parse should fail
	//

	url := "://foobar.com"
	_, err := Get(url, nil, nil, nil)
	assert.NotEqual(t, nil, err)
	//
	// Unsupported protocol scheme - HttpClient.Do should fail
	//
	url = "foo://bar.com"
	_, err = Get(url, nil, nil, nil)
	assert.NotEqual(t, nil, err)
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGetBasicAuth))
	defer srv.Close()
	s := Session{}
	s.UnsafeBasicAuth = true // Otherwise we will get error with httptest
	r := Request{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: "GET",
		Opts: &Opts{
			Userinfo:       url.UserPassword("jtkirk", "Beam me up, Scotty!"),
			ExpectedStatus: 200,
		},
	}
	_, err := s.Send(&r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnsafeBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePost))
	defer srv.Close()
	r := Request{
		Url:    "http://" + srv.Listener.Addr().String(),
		Method: "GET",
		Opts: &Opts{
			Userinfo: url.UserPassword("a", "b"),
		},
	}
	_, err := Send(&r)
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
	fooParams = Params{"foo": "bar"}
	barParams = Params{"bar": "baz"}
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
	url := "http://" + srv.Listener.Addr().String()
	p := fooParams
	res := structType{}
	resp, err := Get(url, &p, &res, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, resp.Status())
	assert.Equal(t, res, barStruct)
	//
	// Bad request
	//
	url = "http://" + srv.Listener.Addr().String()
	p = Params{"bad": "value"}
	e := errorStruct{}
	opts := Opts{
		ExpectedStatus: 200,
	}
	resp, err = Get(url, &p, nil, &opts)
	if err != UnexpectedStatus {
		t.Error(err)
	}
	assert.Equal(t, 500, resp.Status())
	expected := errorStruct{
		Message: "Bad query params: bad=value",
		Status:  500,
	}
	resp.Unmarshall(&e)
	assert.Equal(t, e, expected)
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePost))
	defer srv.Close()
	s := Session{}
	s.Log = true
	url := "http://" + srv.Listener.Addr().String()
	payload := fooStruct
	res := structType{}
	resp, err := s.Post(url, &payload, &res, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, resp.Status())
	assert.Equal(t, res, barStruct)
}

func TestPostUnmarshallable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePost))
	defer srv.Close()
	type ft func()
	var f ft
	url := "http://" + srv.Listener.Addr().String()
	res := structType{}
	payload := f
	_, err := Post(url, &payload, &res, nil)
	assert.NotEqual(t, nil, err)
	_, ok := err.(*json.UnsupportedTypeError)
	if !ok {
		t.Log(err)
		t.Error("Expected json.UnsupportedTypeError")
	}
}

func TestPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandlePut))
	defer srv.Close()
	url := "http://" + srv.Listener.Addr().String()
	res := structType{}
	resp, err := Put(url, &fooStruct, &res, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Status(), 200)
	// Server should return NO data
	assert.Equal(t, resp.RawText(), "")
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
	for k, _ := range fooParams {
		if fooParams[k] != q.Get(k) {
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

func HandleGetBasicAuth(w http.ResponseWriter, req *http.Request) {
	authRegex := regexp.MustCompile(`[Bb]asic (?P<encoded>\S+)`)
	str := req.Header.Get("Authorization")
	matches := authRegex.FindStringSubmatch(str)
	if len(matches) != 2 {
		msg := "Regex doesn't match"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	encoded := matches[1]
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		msg := "Base64 decode failed"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		msg := "String split failed"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	username := parts[0]
	password := parts[1]
	if username != "jtkirk" || password != "Beam me up, Scotty!" {
		code := http.StatusUnauthorized
		text := http.StatusText(code)
		http.Error(w, text, code)
		return
	}
	w.WriteHeader(200)
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
