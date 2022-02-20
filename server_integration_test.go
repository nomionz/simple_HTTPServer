package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppendingAndRetrieving(t *testing.T) {
	store := NewInMemoryPMStore()
	srv := ProjectManagementServer{store}
	worker := "John"

	srv.ServeHTTP(httptest.NewRecorder(), newPostAppendReq(worker))
	srv.ServeHTTP(httptest.NewRecorder(), newPostAppendReq(worker))

	response := httptest.NewRecorder()
	srv.ServeHTTP(response, newGetTasksRequest(worker))

	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "2")
}
