// Example demonstrating use of package restclient, with HTTP Basic
// authentictation over HTTPS, to retrieve a Github auth token.
package main

/*

NOTE: This example may only work on *nix systems due to gopass requirements.

*/

import (
	"code.google.com/p/gopass"
	"fmt"
	"github.com/jmcvetta/restclient"
	"log"
	"net/url"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	//
	// Prompt user for Github username/password
	//
	var username string
	fmt.Printf("Github username: ")
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		log.Fatal(err)
	}
	passwd, err := gopass.GetPass("Github password: ")
	if err != nil {
		log.Fatal(err)
	}
	//
	// Compose request
	//
	// http://developer.github.com/v3/oauth/#create-a-new-authorization
	//
	d := struct {
		Scopes []string `json:"scopes"`
		Note   string   `json:"note"`
	}{
		Scopes: []string{"public_repo"},
		Note:   "testing Go restclient",
	}
	//
	// Struct to hold response data
	//
	res := struct {
		Id        int
		Url       string
		Scopes    []string
		Token     string
		App       map[string]string
		Note      string
		NoteUrl   string `json:"note_url"`
		UpdatedAt string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
	}{}
	//
	// Struct to hold error response
	//
	e := struct {
		Message string
	}{}
	//
	// Setup HTTP Basic auth (ONLY use this with SSL)
	//
	u := url.UserPassword(username, passwd)
	rr := restclient.RequestResponse{
		Url:      "https://api.github.com/authorizations",
		Userinfo: u,
		Method:   "POST",
		Data:     &d,
		Result:   &res,
		Error:    &e,
	}
	//
	// Send request to server
	//
	status, err := restclient.Do(&rr)
	if err != nil {
		log.Fatal(err)
	}
	//
	// Process response
	//
	println("")
	if status == 201 {
		fmt.Printf("Github auth token: %s\n\n", res.Token)
	} else {
		fmt.Println("Bad response status from Github server")
		fmt.Printf("\t Status:  %v\n", status)
		fmt.Printf("\t Message: %v\n", e.Message)
	}
	println("")
}
