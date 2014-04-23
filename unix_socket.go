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
func (d SocketTransport) RoundTrip(req *http.Request) (*http.Response, error) {
  dial, dialErr := net.Dial("unix", d.path)
  if dialErr != nil {
    return nil, dialErr
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
  return "", "", errors.New("No Unix socket found in " + rawPath)
}
