package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Emojis struct {
	Username string `json:"user_name"`
	Emoji    string `json:"emoji"`
}

var DB *sql.DB
var err error

func init() {
	DB, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tl_ancillaries")

	if err != nil {
		panic(err.Error())
	}
}

func emojiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		URLQuery := r.URL.Query()
		user_name := URLQuery["user_name"][0]
		result, err := DB.Query("SELECT user_name, emoji from emojis WHERE user_name = ?", user_name)
		if err != nil {
			panic(err.Error())
		}
		var emoji Emojis

		for result.Next() {
			err = result.Scan(&emoji.Username, &emoji.Emoji)

			if err != nil {
				panic(err.Error())
			}
		}
		fmt.Fprint(w, emoji)

	case "POST":
		d := json.NewDecoder(r.Body)
		emoji := &Emojis{}
		err := d.Decode(emoji)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		DB.Query("INSERT INTO emojis(user_name, emoji) VALUES (?, ?);", emoji.Username, emoji.Emoji)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")

	}
}

func main() {
	http.HandleFunc("/emoji", emojiHandler)

	log.Println("Go!")
	http.ListenAndServe(":8001", nil)

	defer DB.Close()

}

/*
	var httpQueries = ""
	rawQuery := r.URL.Query()
	emojis := rawQuery["text"][0]
	httpQueries = httpQueries + emojis
	fmt.Fprintf(w, httpQueries)
*/
