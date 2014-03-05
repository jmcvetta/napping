package httpclient

import (
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

var starter sync.Once
var addr net.Addr

func testHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(500 * time.Millisecond)
	io.WriteString(w, "hello, world!\n")
}

func testDelayedHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(2100 * time.Millisecond)
	io.WriteString(w, "hello, world ... in a bit\n")
}

func setupMockServer(t *testing.T) {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/test-delayed", testDelayedHandler)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen - %s", err.Error())
	}
	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			t.Fatalf("failed to start HTTP server - %s", err.Error())
		}
	}()
	addr = ln.Addr()
}

func TestHttpClient(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })

	req, _ := http.NewRequest("GET", "http://"+addr.String()+"/test", nil)

	connectTimeout := (250 * time.Millisecond)
	readWriteTimeout := (50 * time.Millisecond)

	httpClient := NewWithTimeout(connectTimeout, readWriteTimeout)

	resp, err := httpClient.Do(req)
	if err == nil {
		t.Fatalf("2nd request should have timed out")
	}

	resp, err = httpClient.Do(req)
	if resp != nil {
		t.Fatalf("3nd request should not have timed out")
	}

}
