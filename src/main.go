package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/database"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/cards", getCards)

	port := ":8090"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func ping(w http.ResponseWriter, req *http.Request) {
	dbcs := database.GetDatabaseConnectionString()
	err := database.Ping(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}
	fmt.Fprintf(w, "successfully connected to database\n")
}

func getCards(w http.ResponseWriter, req *http.Request) {
	dbcs := database.GetDatabaseConnectionString()
	cards, err := database.GetJudgeCards(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	tmpl, err := template.ParseFiles("templates/cards.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", cards)
}
