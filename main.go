package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

func (c *Comment) MemDatConvert() string {
	return "{\n\tID: " + c.ID + ",\n\tBody: " + c.Body + ",\n\tReplies: [\n\t\t" + strings.Join(c.Replies[:], ",\n\t\t") + "\n\t]\n}"
}

func main() {
	fmt.Print((&Comment{CommentStem{"ee", "no u"}, []string{"hello", "fuck you"}}).MemDatConvert())

	database := *newDatabase()
	database.add(Comment{CommentStem{"ee", "no u"}, []string{"hello", "fuck you"}})
	fmt.Println(database.Documents[0])

	fmt.Println(func(i interface{}) interface{} { return i }(Comment{CommentStem{"ee", "no u"}, []string{"hello", "fuck you"}}))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "EE")
		fmt.Println(req.Header)
	})
	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { mux.ServeHTTP(w, req) })))
}
