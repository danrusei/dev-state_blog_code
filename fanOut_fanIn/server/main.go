package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var first sync.Once
var comments []Comment

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/comments/{article:[0-9]+}", CommentsHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//HomeHandler explain the scope of this server
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "This server is used for testing, it serves an json")
}

//CommentsHandler servers the article comments
func CommentsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	first.Do(func() {
		comments = loadJSON("temp.json")
	})

	postID := vars["article"]
	id, err := strconv.Atoi(postID)
	if err != nil {
		log.Fatal("Can't convert teh postID from string to int")
	}

	results := []Comment{}

	for _, comment := range comments {
		if comment.PostID == id {
			results = append(results, comment)
		}

	}

	s, err := json.Marshal(results)
	if err != nil {
		log.Fatal("Can't marshall the result to bytes")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(s)
}

//Comment holds the article comments
type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func loadJSON(filename string) []Comment {

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Can't open the file")
	}

	comments := []Comment{}
	if err := json.Unmarshal(fileContent, &comments); err != nil {
		log.Fatal("Can't unmarshal the contents of the file")
	}

	return comments
}
