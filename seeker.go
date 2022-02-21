package main

import (
	"os"
)

// seeker is just a solution for situation if I'll want to implement delete
type seeker struct {
	file *os.File
}

func (s *seeker) Write(w []byte) (int, error) {
	s.file.Truncate(0)
	s.file.Seek(0, 0)
	return s.file.Write(w)
}
