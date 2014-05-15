package main

import (
  "fmt"
  "github.com/jmcvetta/napping"
)

func main() {
  s := napping.Session{Log: true}
  // Hit the Docker remote API route (v1.10) that shows the info for all local images.
  _, err := s.Get("unix://var/run/docker.sock/v1.10/images/json", nil, nil, nil)
  if err != nil {
    fmt.Println(err)
  }
}