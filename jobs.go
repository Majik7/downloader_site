package main

import "sync"

type jobStatus string

const (
	statusPending jobStatus = "pending"
	statusDone    jobStatus = "done"
	statusError   jobStatus = "error"
)

type job struct {
	Status jobStatus
	Path   string
	Err    string
}

type jobStore struct {
	mu   sync.RWMutex
	jobs map[string]*job
}

func newJobStore() *jobStore {
	return &jobStore{jobs: make(map[string]*job)}
}

func (s *jobStore) create(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[id] = &job{Status: statusPending}
}

func (s *jobStore) get(id string) (*job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	return j, ok
}

func (s *jobStore) setDone(id, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		j.Status = statusDone
		j.Path = path
	}
}

func (s *jobStore) setError(id, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		j.Status = statusError
		j.Err = msg
	}
}

func (s *jobStore) delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
}
