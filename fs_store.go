package main

import (
	"encoding/json"
	"io"
)

// FSStore is a simple file system storage
type FSStore struct {
	file io.ReadWriteSeeker
}

func (f *FSStore) GetProjectInfo() Project {
	// want to be able to read twice
	f.file.Seek(0, 0)
	projectInfo, _ := NewProject(f.file)
	return projectInfo
}

func (f *FSStore) GetDoneTasks(name string) int {
	worker := f.GetProjectInfo().Find(name)
	if worker != nil {
		return worker.DoneTasks
	}
	return 0
}

func (f *FSStore) Append(name string) {
	project := f.GetProjectInfo()
	worker := project.Find(name)
	if worker != nil {
		worker.DoneTasks++
	} else {
		project = append(project, Worker{name, 1})
	}
	f.file.Seek(0, 0)
	json.NewEncoder(f.file).Encode(&project)
}
