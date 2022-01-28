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
		temp, err := template.ParseFiles("templates/adventure.html")
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

type homeHandler struct {
	adv advHandler
}

func (h homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "" || path == "home" {
		home, err := os.ReadFile("templates/home.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, string(home))
		}
	} else {
		h.adv.ServeHTTP(w, r)
	}
}

type Game interface {
	start(m map[string]adventure) error
}

type WebGame struct{}

func (wg WebGame) start(m map[string]adventure) error {
	http.Handle("/", homeHandler{adv: advHandler{adventures: m}})
	return http.ListenAndServe(":8080", nil)
}

type CLIGame struct{}

func (cliG CLIGame) start(m map[string]adventure) error {
	fmt.Println("Start game in cli mode")

	nextVenture := "intro"

	for nextVenture != "home" {
		venture := m[nextVenture]
		fmt.Println(venture.Story)

		l := len(venture.Options)
		for i := 0; i < l; i++ {
			fmt.Println(venture.Options[i].Text)
			fmt.Println("Press", i+1, "to venture to", venture.Options[i].Arc)
		}

		var choice int
		for choice < 1 || choice > 2 {
			_, err := fmt.Scanf("%d", &choice)
			if err != nil {
				log.Fatal(err)
			}
		}

		nextVenture = venture.Options[choice-1].Arc
	}

	os.Exit(0)

	return nil
}

func newGame(mode string) Game {
	switch mode {
	case "web":
		return WebGame{}
	case "cli":
		return CLIGame{}
	default:
		log.Fatal("Unknown game mode")
		return nil
	}
}

func loadAdventures(filename string) map[string]adventure {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	m := map[string]adventure{}
	err = json.Unmarshal(f, &m)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func main() {
	mode := flag.String("i", "cli", "Interface to play the game, web or cli.")
	advFilename := flag.String("f", "adventure.json", "Path of the file containing the adventures.")
	flag.Parse()

	m := loadAdventures(*advFilename)

	game := newGame(*mode)
	err := game.start(m)
	if err != nil {
		log.Fatal(err)
	}
}
