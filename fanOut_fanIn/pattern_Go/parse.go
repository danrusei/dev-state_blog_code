package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type post struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func parse(done chan struct{}, id int) string {

	body, err := getBody(id)
	//close done when we get an error while parsing
	if err != nil {
		close(done)
		log.Fatal(err)
	}

	posts := []post{}

	if err := json.Unmarshal(body, &posts); err != nil {
		close(done)
		log.Fatalf("Can't unmarshal: %v", err)
	}

	longestPost := 0
	longestPostID := 0
	longestPostEmail := ""

	for _, p := range posts {
		if len(p.Body) > longestPost {
			longestPost = len(p.Body)
			longestPostID = p.PostID
			longestPostEmail = p.Email
		}
	}

	return fmt.Sprintf("%d %s %d", longestPostID, longestPostEmail, longestPost)

}

func getBody(id int) ([]byte, error) {
	site, err := url.Parse("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		return nil, err
	}
	q := site.Query()
	q.Set("postId", strconv.Itoa(id))
	site.RawQuery = q.Encode()
	log.Println("Getting: ", site.String())

	client := &http.Client{}
	req, err := http.NewRequest("GET", site.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
