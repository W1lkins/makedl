package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const makefileURL = "https://github.com/evalexpr/makefiles"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "makedl",
	Short: fmt.Sprintf("Download Makefiles from %s", makefileURL),
	Long: fmt.Sprintf(`Download Makefiles from %s.

This application is a tool to fetch Makefiles from GitHub for a range of
languages.`, makefileURL),
	Run: Download,
}

func Download(cmd *cobra.Command, args []string) {
	err := download(args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func download(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: makedl [lang]")
	}

	if _, err := os.Stat("./Makefile"); !os.IsNotExist(err) {
		return errors.New("refusing to write Makefile, already exists")
	}

	language := args[0]
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

	url := fmt.Sprintf("https://raw.githubusercontent.com/evalexpr/makefiles/master/%s.Makefile", language)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("makefile not found for language: %s", language)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("writing new makefile for language: %s", language)
	if err := ioutil.WriteFile("./Makefile", bytes, 0644); err != nil {
		return fmt.Errorf("could not write Makefile: %w", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
