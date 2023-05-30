package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var tom *person = &person{
	Name: "Tom",
	Age:  28,
}

var httpQueries = ""

func emojiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rawQuery := r.URL.Query()
		emojis := rawQuery["text"][0]
		httpQueries = httpQueries + emojis
		fmt.Fprintf(w, httpQueries)
	case "POST":
		d := json.NewDecoder(r.Body)
		p := &person{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tom = p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")

	}
}

func main() {
	http.HandleFunc("/emoji", emojiHandler)

	log.Println("Go!")
	http.ListenAndServe(":8001", nil)
}
