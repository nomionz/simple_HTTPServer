package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppendingAndRetrieving(t *testing.T) {
	file, rmFile := createTmpFile(t, `[]`)
	defer rmFile()
	store, err := NewFSStore(file)
	assertNoError(t, err)
	srv := NewPMServer(store)
	worker := "John"

	srv.ServeHTTP(httptest.NewRecorder(), newPostAppendReq(worker))
	srv.ServeHTTP(httptest.NewRecorder(), newPostAppendReq(worker))

	t.Run("get num of done tasks", func(t *testing.T) {
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, newGetTasksRequest(worker))

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "2")
	})

	t.Run("get project info", func(t *testing.T) {
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, newProjectRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getProjectInfoFromResponse(t, response.Body)
		want := []Worker{{"John", 2}}
		assertProject(t, got, want)
	})
}
