package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Worker is the struct for JSON
type Worker struct {
	Name      string
	DoneTasks int
}

// ProjectManagementStore is the interface for PMServer
type ProjectManagementStore interface {
	GetDoneTasks(name string) int
	Append(name string)
	GetProjectInfo() []Worker
}

type ProjectManagementServer struct {
	store ProjectManagementStore
	http.Handler
}

// NewPMServer takes ProjectManagementStore as a parameter, set up router and return instance of a server
func NewPMServer(store ProjectManagementStore) *ProjectManagementServer {
	p := new(ProjectManagementServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/project", http.HandlerFunc(p.projectHandler))
	router.Handle("/workers/", http.HandlerFunc(p.workersHandler))
	p.Handler = router
	return p
}

// projectHandler handles /project endpoint
func (p *ProjectManagementServer) projectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetProjectInfo())
}

// workersHandler handles /workers/ endpoint
func (p *ProjectManagementServer) workersHandler(w http.ResponseWriter, r *http.Request) {
	worker := strings.TrimPrefix(r.URL.Path, "/workers/")
	switch r.Method {
	case http.MethodPost:
		p.processAppend(w, worker)
	case http.MethodGet:
		p.showTasks(w, worker)
	}
}

// showTasks handler of a GET http request
func (p *ProjectManagementServer) showTasks(w http.ResponseWriter, worker string) {
	tasks := p.store.GetDoneTasks(worker)
	if tasks == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, tasks)
}

// processAppend handler of a POST http request
func (p *ProjectManagementServer) processAppend(w http.ResponseWriter, worker string) {
	p.store.Append(worker)
	w.WriteHeader(http.StatusAccepted)
}
