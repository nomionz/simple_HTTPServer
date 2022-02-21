package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Spies for mocking
type StubPMStore struct {
	tasks map[string]int
	// Detect POST calls
	doneCalls []string
}

func (s *StubPMStore) GetDoneTasks(name string) int {
	return s.tasks[name]
}

func (s *StubPMStore) Append(name string) {
	s.doneCalls = append(s.doneCalls, name)
}

func TestGETPizzas(t *testing.T) {
	store := StubPMStore{
		map[string]int{
			"John":  20,
			"Steve": 10,
		},
		nil,
	}
	server := NewPMServer(&store)

	t.Run("returns John's done tasks", func(t *testing.T) {
		request := newGetTasksRequest("John")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Steve's done tasks", func(t *testing.T) {
		request := newGetTasksRequest("Steve")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing pizzas", func(t *testing.T) {
		request := newGetTasksRequest("Benji")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreDone(t *testing.T) {
	store := StubPMStore{
		map[string]int{},
		nil,
	}
	server := NewPMServer(&store)

	t.Run("it records wins on POST", func(t *testing.T) {
		worker := "John"

		request := newPostAppendReq(worker)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.doneCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(store.doneCalls), 1)
		}

		if store.doneCalls[0] != worker {
			t.Errorf("did not store correct winner got %q want %q", store.doneCalls[0], worker)
		}
	})
}

//server_test.go
func TestLeague(t *testing.T) {
	store := StubPMStore{}
	server := NewPMServer(&store)

	t.Run("it returns 200 on /project", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/project", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func newPostAppendReq(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/workers/%s", name), nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetTasksRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/workers/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
