package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Worker struct {
	Name      string
	DoneTasks int
}

type ProjectManagementStore interface {
	GetDoneTasks(name string) int
	Append(name string)
	GetProjectInfo() []Worker
}

type ProjectManagementServer struct {
	store ProjectManagementStore
	http.Handler
}

func NewPMServer(store ProjectManagementStore) *ProjectManagementServer {
	p := new(ProjectManagementServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/project", http.HandlerFunc(p.projectHandler))
	router.Handle("/workers/", http.HandlerFunc(p.workersHandler))
	p.Handler = router
	return p
}

func (p *ProjectManagementServer) projectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.GetProjectInfo())
}

func (p *ProjectManagementServer) workersHandler(w http.ResponseWriter, r *http.Request) {
	worker := strings.TrimPrefix(r.URL.Path, "/workers/")
	switch r.Method {
	case http.MethodPost:
		p.processAppend(w, worker)
	case http.MethodGet:
		p.showTasks(w, worker)
	}
}

func (p *ProjectManagementServer) showTasks(w http.ResponseWriter, worker string) {
	tasks := p.store.GetDoneTasks(worker)
	if tasks == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, tasks)
}

func (p *ProjectManagementServer) processAppend(w http.ResponseWriter, worker string) {
	p.store.Append(worker)
	w.WriteHeader(http.StatusAccepted)
}
