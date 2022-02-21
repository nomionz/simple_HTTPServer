package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type Project []Worker

func NewProject(r io.Reader) (Project, error) {
	var project Project
	err := json.NewDecoder(r).Decode(&project)
	if err != nil {
		err = fmt.Errorf("Couldn't parse project info '%v'", err)
	}
	return project, err
}

func (p Project) Find(name string) *Worker {
	for i, worker := range p {
		if worker.Name == name {
			return &p[i]
		}
	}
	return nil
}
