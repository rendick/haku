package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

var engines = map[int]string{
	0: "https://lite.duckduckgo.com/lite/?q=",
}

const choice = 0

var url string = ""
var formattedInput string

func Home(w http.ResponseWriter, req *http.Request) {
	var title []string
	var links []string
	var description []string

	for k, _ := range engines {
		if choice == k {
			url = engines[k]
		}
	}

	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Connecting ", r.URL)
	})

	url = url + formattedInput

	collector.OnHTML("a.result-link", func(h *colly.HTMLElement) {
		title = append(title, h.Text)
		links = append(links, h.Attr("href"))
	})
	collector.OnHTML("td.result-snippet", func(h *colly.HTMLElement) {
		description = append(description, h.Text)
	})

	collector.Visit(url)

	tmplt, err := template.ParseFiles("client/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title       []string
		Links       []string
		Description []string
		Entered     string
	}{
		Title:       title,
		Links:       links,
		Description: description,
		Entered:     strings.ReplaceAll(formattedInput, "%20", " "),
	}
	tmplt.Execute(w, data)
}

func Settings(w http.ResponseWriter, req *http.Request) {
	tmplt, err := template.ParseFiles("client/settings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplt.Execute(w, nil)
}

func Server() {
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))

	http.HandleFunc("/", Home)
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		handleInputText := r.FormValue("user-input")
		formattedInput = strings.ReplaceAll(handleInputText, " ", "%20")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	http.HandleFunc("/settings", Settings)

	http.ListenAndServe(":8080", nil)
}
