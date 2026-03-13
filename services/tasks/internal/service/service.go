package service

import (
	"fmt"
	"sync"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type TaskService struct {
	mu    sync.RWMutex
	tasks map[string]*Task
	seq   int
}

func New() *TaskService {
	return &TaskService{tasks: make(map[string]*Task)}
}

func (s *TaskService) Create(title, description, dueDate string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	t := &Task{
		ID:          fmt.Sprintf("t_%03d", s.seq),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Done:        false,
	}
	s.tasks[t.ID] = t
	return t
}

func (s *TaskService) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

func (s *TaskService) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskService) Update(id, title string, done *bool, description string) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	if title != "" {
		t.Title = title
	}
	if done != nil {
		t.Done = *done
	}
	if description != "" {
		t.Description = description
	}
	return t, true
}

func (s *TaskService) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if ok {
		delete(s.tasks, id)
	}
	return ok
}
