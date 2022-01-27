package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type adventure struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type advHandler struct {
	adventures map[string]adventure
}

func (h advHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	adv := r.URL.Path[1:]
	if v, ok := h.adventures[adv]; ok {
		temp, err := template.ParseFiles("adventure.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}

		err = temp.Execute(w, v)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	advFilename := flag.String("f", "adventure.json", "Path of the file containing the adventures.")

	m, err := loadAdventures(*advFilename)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", advHandler{adventures: m})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadAdventures(filename string) (map[string]adventure, error) {
	advJson, err := os.ReadFile("adventure.json")
	if err != nil {
		return nil, err
	}

	m := map[string]adventure{}
	err = json.Unmarshal(advJson, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
