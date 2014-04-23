package napping

import (
  "net/url"
  "net"
  "net/http"
  "net/http/httputil"
  "errors"
  "path"
  "os/exec"
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
  cmd := exec.Command("file", "-b", path)
  out, err := cmd.Output()
  if err != nil {
      return false
  }
  return string(out) == "socket\n"
}

func LocateSocket(rawUrl string) (string, string, error) {
  u, err := url.Parse(rawUrl)
  if err != nil {
    return "", "", err
  }
  if u.Scheme != "unix" {
    return "", "", errors.New("URL scheme must be 'unix'")
  }
  s := u.Path
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
  return "", "", errors.New("No Unix socket found in " + rawUrl)
}
