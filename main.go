package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Things I need to make this work:

/*
- an agent that is instantiated with the following at minimum:
  - a service (https://bsky.social)
  - a login method
  - USE OAUTH
  - a post method
*/

func main() {
	username := os.Getenv("BLUESKY_HANDLE")
	appPassword := os.Getenv("BLUESKY_PASSWORD")
	agent := BskyAgent{
		client: &http.Client{},
		host:   "https://bsky.social",
	}
	agent.login(username, appPassword)
}

type BskyAgent struct {
	client  *http.Client
	session *SessionToken
	host    string
}

type SessionToken struct {
	Did    string `json:"did"`
	DidDoc struct {
		Context            []string `json:"@context"`
		ID                 string   `json:"id"`
		AlsoKnownAs        []string `json:"alsoKnownAs"`
		VerificationMethod []struct {
			ID                 string `json:"id"`
			Type               string `json:"type"`
			Controller         string `json:"controller"`
			PublicKeyMultibase string `json:"publicKeyMultibase"`
		} `json:"verificationMethod"`
		Service []struct {
			ID              string `json:"id"`
			Type            string `json:"type"`
			ServiceEndpoint string `json:"serviceEndpoint"`
		} `json:"service"`
	} `json:"didDoc"`
	Handle          string `json:"handle"`
	Email           string `json:"email"`
	EmailConfirmed  bool   `json:"emailConfirmed"`
	EmailAuthFactor bool   `json:"emailAuthFactor"`
	AccessJwt       string `json:"accessJwt"`
	RefreshJwt      string `json:"refreshJwt"`
	Active          bool   `json:"active"`
}

type Post struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Record     struct {
		Text      string `json:"text"`
		CreatedAt string `json:"createdAt"`
	} `json:"record"`
	Response struct {
		URI string `json:"uri"`
		CID string `json:"cid"`
	}
}

func (b BskyAgent) login(u, p string) {
}

func (b BskyAgent) post(c string) {
	body := strings.NewReader("test post")
	b.client.Post("test", "application/json", body)
}
