package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("project from a reader", func(t *testing.T) {
		file, rmFile := createTmpFile(t, `[
            {"Name": "John", "DoneTasks": 10},
            {"Name": "Chris", "DoneTasks": 33}]`)
		defer rmFile()

		store, err := NewFSStore(file)
		assertNoError(t, err)

		got := store.GetProjectInfo()

		want := []Worker{
			{"John", 10},
			{"Chris", 33},
		}
		// want to be able to read twice
		got = store.GetProjectInfo()
		assertProject(t, got, want)
	})

	t.Run("Get worker's done tasks", func(t *testing.T) {
		file, rmFile := createTmpFile(t, `[
            {"Name": "John", "DoneTasks": 10},
            {"Name": "Chris", "DoneTasks": 33}]`)
		defer rmFile()

		store, err := NewFSStore(file)
		assertNoError(t, err)

		got := store.GetDoneTasks("John")
		want := 10
		assertTasksEquals(t, got, want)

	})

	t.Run("record done task for an existing worker", func(t *testing.T) {
		file, rmFile := createTmpFile(t, `[
            {"Name": "John", "DoneTasks": 10},
            {"Name": "Chris", "DoneTasks": 33}]`)
		defer rmFile()

		store, err := NewFSStore(file)
		assertNoError(t, err)
		store.Append("John")

		got := store.GetDoneTasks("John")
		want := 11
		assertTasksEquals(t, got, want)
	})

	t.Run("store done task of new worker", func(t *testing.T) {
		file, rmFile := createTmpFile(t, `[
            {"Name": "John", "DoneTasks": 10},
            {"Name": "Chris", "DoneTasks": 33}]`)
		defer rmFile()
		store, err := NewFSStore(file)
		assertNoError(t, err)
		store.Append("Steve")

		got := store.GetDoneTasks("Steve")
		want := 1
		assertTasksEquals(t, got, want)

	})

	t.Run("works with an empty file", func(t *testing.T) {
		file, rmFile := createTmpFile(t, "")
		defer rmFile()

		_, err := NewFSStore(file)

		assertNoError(t, err)
	})

}

func assertTasksEquals(t testing.TB, got, want int) {
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func createTmpFile(t testing.TB, data string) (*os.File, func()) {
	t.Helper()

	tmp, err := ioutil.TempFile("", "tmpFileForTest")

	if err != nil {
		t.Fatalf("Couldn't create tmp file '%v'", err)
	}

	tmp.Write([]byte(data))
	// I want to remove file after tested it out so no memory leaks during the simple test
	rmFile := func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}

	return tmp, rmFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one %v", err)
	}
}
