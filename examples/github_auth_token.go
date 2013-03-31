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
	"github.com/kr/pretty"
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
	// Setup HTTP Basic auth (ONLY use this with SSL)
	//
	u := url.UserPassword("jmcvetta", passwd)
	rr := restclient.RequestResponse{
		Url:      "https://api.github.com/authorizations",
		Userinfo: u,
		Method:   "POST",
		Data:     &d,
		Result:   &res,
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
		return
	}
	fmt.Println("Bad response from Github server:")
	pretty.Printf("%# v\n", res)
}
