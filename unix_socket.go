package napping

import (
  "net"
  "net/http"
  "net/http/httputil"
  "errors"
  "path"
  "os"
  "fmt"
)

type SocketTransport struct {path string}

// The RoundTripper (http://golang.org/pkg/net/http/#RoundTripper) for the socket transport dials the socket
// each time a request is made.
func (d SocketTransport) RoundTrip(req *http.Request) (*http.Response, error) {
  dial, err := net.Dial("unix", d.path)
  if err != nil {
    return nil, err
  }
  socketClientConn := httputil.NewClientConn(dial, nil)
  defer socketClientConn.Close()
  return socketClientConn.Do(req)
}

func isUnixSocket(path string) bool {
  fi, err := os.Lstat(path)
  if err != nil {
    fmt.Println(path + " " + err.Error())
    return false
  }
  return fi.Mode()&os.ModeType == os.ModeSocket
}

// Unix urls like unix://var/run/docker.sock/v1.10/images/json contain two parts: the path to the Unix domain
// socket, /var/run/docker.sock, and the request path: /v1.10/images/json . This function splits out the two
// parts by walking down the path to / and checking whether each node is a Unix domain socket as opposed to
// regular file or directory.
//
// If no Unix domain sockets are found, an error is returned.
func LocateSocket(rawPath string) (string, string, error) {
  s := rawPath
  if s[0] != '/' {
    s = "/" + s
  }
  p := ""
  p_ := ""
  for s != "" {
    if last := len(s) - 1; last >= 0 && s[last] == '/' {
      s = s[:last]
    }
    if isUnixSocket(s) {
      return s, "/" + p, nil
    }
    s, p_ = path.Split(s)
    p = path.Join(p_, p)
  }
  return "", "", errors.New("No Unix domain socket found in " + rawPath)
}
