package main

import (
	"io/ioutil"
	"testing"
)

func TestSeeker(t *testing.T) {
	file, rmFile := createTmpFile(t, "1234")
	defer rmFile()

	tape := &seeker{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
