package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Spies for mocking
type StubPMStore struct {
	tasks map[string]int
	// Detect POST calls
	doneCalls []string
	project   Project
}

func (s *StubPMStore) GetDoneTasks(name string) int {
	return s.tasks[name]
}

func (s *StubPMStore) Append(name string) {
	s.doneCalls = append(s.doneCalls, name)
}

func (s *StubPMStore) GetProjectInfo() Project {
	return s.project
}

func TestGETPizzas(t *testing.T) {
	store := StubPMStore{
		map[string]int{
			"John":  20,
			"Steve": 10,
		},
		nil,
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

	t.Run("returns 404 on missing tasks", func(t *testing.T) {
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

func TestLeague(t *testing.T) {

	t.Run("it returns JSON on /project", func(t *testing.T) {
		want := Project{
			{"John", 10},
			{"Steve", 20},
			{"Martin", 13},
		}
		store := StubPMStore{nil, nil, want}
		srv := NewPMServer(&store)

		request := newProjectRequest()
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, request)

		got := getProjectInfoFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertProject(t, got, want)
		if response.Result().Header.Get("content-type") != "application/json" {
			t.Errorf("response didn't have json got %v", response.Result().Header)
		}
	})
}

func assertProject(t testing.TB, got, want Project) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func getProjectInfoFromResponse(t testing.TB, body io.Reader) (project Project) {
	err := json.NewDecoder(body).Decode(&project)

	if err != nil {
		t.Fatalf("Couldn't parse response from %q into Worker's slice '%v'", body, err)
	}
	return
}

func newProjectRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/project", nil)
	return req
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
