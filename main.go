package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Grabbing login details...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("BLUESKY_HANDLE")
	appPassword := os.Getenv("BLUESKY_PASSWORD")
	fmt.Println("Creating new session...")
	agent := &BskyAgent{
		Client:   http.Client{},
		Host:     "https://bsky.social",
		Session:  SessionToken{},
		Posts:    []Posts{},
		IoReader: bufio.NewScanner(os.Stdin),
	}
	agent.login(username, appPassword)
	fmt.Println("Connected!")
	fmt.Print("Please enter your post:\n> ")
	post := agent.getInput()
	fmt.Println("Posting...")
	agent.post(post)
	fmt.Println("Shutting down...")
	os.Exit(0)
}

type BskyAgent struct {
	Client   http.Client
	Session  SessionToken
	Host     string
	Posts    []Posts
	IoReader *bufio.Scanner
}

type SessionToken struct {
	DID    string `json:"did"`
	DIDDoc struct {
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
	AccessJWT       string `json:"accessJwt"`
	RefreshJWT      string `json:"refreshJwt"`
	Active          bool   `json:"active"`
}

type Posts struct {
	Post     Post
	Response PostResponse
}
type Post struct {
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

type PostRecord struct {
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type PostResponse struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

func (b *BskyAgent) login(u, p string) {
	login := struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{
		Identifier: u,
		Password:   p,
	}
	payload, err := json.Marshal(login)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payload)
	fmt.Println("trying to login....")
	resp, err := b.Client.Post(b.Host+"/xrpc/com.atproto.server.createSession", "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	result, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		for k, v := range resp.Header {
			fmt.Printf("Key: %s, Value: %s\n", k, v)
		}
		log.Fatalf("Response failed with status code: %d and\nbody %s\n", resp.StatusCode, result)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(result, &b.Session)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Logged in as %s!\n", b.Session.Handle)
}

func (b *BskyAgent) post(c string) {
	createdAt := time.Now().Format(time.RFC3339)
	post := Post{
		Repo:       b.Session.Handle,
		Collection: "app.bsky.feed.post",
		Record: PostRecord{
			Text:      c,
			CreatedAt: createdAt,
		},
	}
	payload, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", b.Host+"/xrpc/com.atproto.repo.createRecord", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.Session.AccessJWT))
	req.Header.Add("Content-Type", "application/json")
	resp, err := b.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	result, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		for k, v := range resp.Header {
			fmt.Printf("Key: %s, Value: %s\n", k, v)
		}
		log.Fatalf("Response failed with status code: %d and\nbody %s\n", resp.StatusCode, result)
	}
	if err != nil {
		log.Fatal(err)
	}
	var response PostResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		log.Fatal(err)
	}
	b.Posts = append(b.Posts, Posts{
		Post:     post,
		Response: response,
	})
	fmt.Printf("Posted %s at %s!\n", c, createdAt)
}

func (b *BskyAgent) getInput() string {
	b.IoReader.Scan()
	return b.IoReader.Text()
}
