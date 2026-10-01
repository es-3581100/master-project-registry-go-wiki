package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store struct {
	path     string
	mu       sync.RWMutex
	projects []Project
}

func LoadConfig(path string) (Config, error) {
	var c Config
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	if c.Bind == "" {
		c.Bind = "127.0.0.1:8788"
	}
	if c.SiteTitle == "" {
		c.SiteTitle = "Master Project Registry"
	}
	return c, nil
}

func OpenStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) reload() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var p []Project
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	SortProjects(p)
	s.projects = p
	return nil
}

func (s *Store) All() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Project, len(s.projects))
	copy(out, s.projects)
	return out
}

func (s *Store) Get(slug string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.projects {
		if p.Slug == slug {
			return p, true
		}
	}
	return Project{}, false
}

func (s *Store) Upsert(p Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.Title = strings.TrimSpace(p.Title)
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Slug == "" {
		p.Slug = Slugify(p.Title)
	}
	if p.ID == "" {
		p.ID = "local:" + p.Slug
	}
	if p.Status == "" {
		p.Status = "active"
	}
	if p.Family == "" {
		p.Family = "Experiments / Other"
	}
	found := false
	for i := range s.projects {
		if s.projects[i].Slug == p.Slug || s.projects[i].ID == p.ID {
			s.projects[i] = p
			found = true
			break
		}
	}
	if !found {
		s.projects = append(s.projects, p)
	}
	SortProjects(s.projects)
	return s.saveLocked()
}

func (s *Store) Mutate(slug string, fn func(*Project)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.projects {
		if s.projects[i].Slug == slug {
			fn(&s.projects[i])
			return s.saveLocked()
		}
	}
	return fmt.Errorf("project %q not found", slug)
}

func (s *Store) saveLocked() error {
	b, err := json.MarshalIndent(s.projects, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
