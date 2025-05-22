package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type LocalStorage struct {
	Engine string `json:"engine"`
}

var engines = map[int]string{
	1: "https://lite.duckduckgo.com/lite/?q=",
	2: "https://www.startpage.com/do/search?query=",
}

var selectedOption int = 0

var url string = ""
var formattedInput string

var (
	title       []string
	links       []string
	description []string
)

func IdentifySearchEngine() {
	for k := range engines {
		if selectedOption == k {
			url = engines[k]
		}
	}
}

func CollectPageInformation() {
	title = nil
	links = nil
	description = nil

	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Connecting", r.URL)
	})

	url = url + formattedInput

	elements := func(titleElement string, linkElement string, descriptionElement string) {
		collector.OnHTML(titleElement, func(h *colly.HTMLElement) {
			title = append(title, h.Text)
		})

		collector.OnHTML(linkElement, func(h *colly.HTMLElement) {
			links = append(links, h.Attr("href"))
		})

		collector.OnHTML(descriptionElement, func(h *colly.HTMLElement) {
			description = append(description, h.Text)
		})
	}

	switch selectedOption {
	case 1:
		elements("a.result-link", "a.result-link", "td.result-snippet")
	case 2:
		elements("h2.wgl-title", "a.result-title", "p.description")
	}

	collector.Visit(url)
}

func Home(w http.ResponseWriter, req *http.Request) {
	IdentifySearchEngine()
	CollectPageInformation()

	tmplt, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(description) > 0 && selectedOption == 2 {
		description = slices.Delete(description, 0, 2)
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
	tmplt, err := template.ParseFiles("web/settings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplt.Execute(w, nil)
}

func CreatePage(name string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(name, handler)
	fmt.Printf("%sINFO:%s Starting %s%s%s page\n", "\033[33;1m", "\033[0m", "\033[;1m", name, "\033[0m")
}

func Server() {
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))

	CreatePage("/json", func(w http.ResponseWriter, r *http.Request) {
		var ls LocalStorage

		err := json.NewDecoder(r.Body).Decode(&ls)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		engine_number, err := strconv.Atoi(ls.Engine)
		if err != nil {
			panic(err)
		}

		selectedOption = engine_number
	})
	CreatePage("/", Home)
	CreatePage("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		handleInputText := r.FormValue("user-input")
		formattedInput = strings.ReplaceAll(handleInputText, " ", "%20")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	CreatePage("/settings", Settings)

	http.ListenAndServe(":8080", nil)
}
