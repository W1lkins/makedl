package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func cleanup(f *os.File) {
	f.Close()
	_ = os.Remove(f.Name())
}

func TestDownload(t *testing.T) {
	var cases = []struct {
		args           []string
		shouldErr      bool
		errMessage     string
		createMakefile bool
	}{
		{
			// No args
			args:           []string{},
			shouldErr:      true,
			errMessage:     "usage: makedl [lang]",
			createMakefile: false,
		},
		{
			// Args for a lang that doesn't exist
			args:           []string{"nope"},
			shouldErr:      true,
			errMessage:     "makefile not found for language: nope",
			createMakefile: false,
		},
		{
			// Makefile already exists
			args:           []string{"go"},
			shouldErr:      true,
			errMessage:     "refusing to write Makefile, already exists",
			createMakefile: true,
		},
	}

	for _, c := range cases {
		if c.createMakefile {
			f, _ := os.Create("./Makefile")
			defer cleanup(f)
		}

		err := download(c.args)
		if err != nil {
			if !c.shouldErr {
				t.Errorf("got unexpected error: %v", err)
			}

			if err.Error() != c.errMessage {
				t.Errorf("got wrong error message, expected: '%s', got '%s'", c.errMessage, err.Error())
			}
		}
	}
}
