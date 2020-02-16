package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const makeFmt = "https://raw.githubusercontent.com/evalexpr/makefiles/master/%s.Makefile"
const path = "./Makefile"

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("usage: %s [language]", os.Args[0])
	}
	language := os.Args[1]

	if language == "" {
		log.Fatalf("lang required in form -lang=foo, -lang bar")
	}

	switch strings.ToLower(language) {
	case "php":
		language = "PHP"
	case "go":
		language = "Go"
	case "golang":
		language = "Go"
	case "rust":
		language = "Rust"
	case "flask":
		language = "Flask"
	}
	log.Printf("looking for makefile for language: %s", language)

	url := fmt.Sprintf(makeFmt, language)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode == http.StatusNotFound {
		log.Fatalf("makefile not found for language: %s", language)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("writing new makefile for language: %s", language)
		ioutil.WriteFile(path, bytes, 0644)
	} else {
		log.Print("refusing to write file since Makefile found in cwd")
	}
}
