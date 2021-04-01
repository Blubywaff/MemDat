package main

import (
	"fmt"
	"log"
	"net/http"
)

type Comment struct {
	CommentStem
	Replies []string `json:"Replies" bson:"Replies"`
}

type CommentStem struct {
	ID   string `json:"ID" bson:"ID"`
	Body string `json:"Body" bson:"Body"`
}

type FullComment struct {
	CommentStem
	Replies []FullComment
}

func main() {
	//fmt.Print(FullComment{"ee", "no u", []string{"hello", "fuck you"}})

	database := newDatabase()
	fmt.Println(database.generateId())

	fmt.Println(func(i interface{}) interface{} { return i }(Comment{CommentStem{"ee", "no u"}, []string{"hello", "fuck you"}}))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "EE")
		fmt.Println(req.Header)
	})
	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { mux.ServeHTTP(w, req) })))
}
