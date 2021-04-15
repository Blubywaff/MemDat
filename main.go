package main

import (
	"fmt"
	"time"
)

type Comment struct {
	CommentStem
	Replies []string `json:"Replies" bson:"Replies" memdat:"Replies"`
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
	/*
		fmt.Println((&Comment{CommentStem{"ee", "no u"}, []string{"hello", "I agree with you"}}).MemDatConvert())
		fmt.Println("e\n\r\te")
		fmt.Println(strings.Contains("e\n\te", "\r"))

		var data interface{}
		str := "{\n\"EE\": [\n\t\"\\\"\\\"\\\"\\n\\\"\\\"\"\n\t]\n}"
		fmt.Println(str)
		fmt.Println(json.Unmarshal([]byte(str), &data))
		fmt.Println(data)
		data = nil
		fmt.Println("{\"EE\": [\"\\\"\\\\\\\"\"]}")
		fmt.Println(json.Unmarshal([]byte("{\"EE\": [\"\\\\\"]}"), &data))
		fmt.Println(data)
	*/

	start := time.Now().UnixNano()

	database := *newDatabase()
	for i := 0; i < 1; i++ {
		database.add(Comment{CommentStem{database.generateId(), database.generateId()}, []string{database.generateId(), database.generateId()}})
	}
	//database.add([]int{1, 2, 3})
	fmt.Println(database.Documents[0])
	database.Documents[0]["new"] = "new"
	fmt.Println(database.Documents[0])

	fmt.Println(time.Now().UnixNano() - start)

	//database.add([]string{"Hello", "21"})

	//fmt.Println(func(i interface{}) interface{} { return i }(Comment{CommentStem{"ee", "no u"}, []string{"hello", "I agree with you"}}))
	/*
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, "EE")
			fmt.Println(req.Header)
		})
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { mux.ServeHTTP(w, req) })))
	*/
}
