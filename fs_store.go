package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// FSStore is a simple file system storage
// pointer to encoder, so I shouldn't call NewEncoder each time I write something to the file
// since I have Project in struct it's more readable in the code
type FSStore struct {
	file    *json.Encoder
	project Project
}

func initFile(file *os.File) error {
	file.Seek(0, 0)
	// Write an empty JSON if file is empty
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("couldn't get file info from '%s' '%v'", file.Name(), err)
	}
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}
	return nil
}

func NewFSStore(file *os.File) (*FSStore, error) {

	err := initFile(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize file '%v'", err)
	}
	proj, err := NewProject(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't load project info store from file %s, %v", file.Name(), err)
	}
	return &FSStore{
		json.NewEncoder(&seeker{file}),
		proj,
	}, nil
}

func (f *FSStore) GetProjectInfo() Project {
	return f.project
}

func (f *FSStore) GetDoneTasks(name string) int {
	worker := f.project.Find(name)
	if worker != nil {
		return worker.DoneTasks
	}
	return 0
}

func (f *FSStore) Append(name string) {
	worker := f.project.Find(name)
	if worker != nil {
		worker.DoneTasks++
	} else {
		f.project = append(f.project, Worker{name, 1})
	}
	f.file.Encode(f.project)
}
