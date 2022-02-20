package main

import (
	"fmt"
	"net/http"
	"strings"
)

type ProjectManagementStore interface {
	GetDoneTasks(name string) int
	Append(name string)
}

type ProjectManagementServer struct {
	store ProjectManagementStore
}

func (p *ProjectManagementServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
