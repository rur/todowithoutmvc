package app

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var CookieName = "todowithoutmvc-session"

type Server interface {
	LoadTodos(*http.Request) (Todos, string)
	SaveTodos(string, Todos) error
}

func NewServer(repo Repository) Server {
	return &server{
		repo: repo,
	}
}

type server struct {
	sync.RWMutex
	repo Repository
}

func (s *server) LoadTodos(req *http.Request) (Todos, string) {
	var key string
	if cookie, err := req.Cookie(CookieName); err == nil {
		key = cookie.Value
	} else {
		return s.repo.NewTodos(), ""
	}
	s.RLock()
	defer s.RUnlock()
	if todos, ok := s.repo[key]; ok {
		return todos, key
	} else {
		return s.repo.NewTodos(), key
	}
}

func (s *server) SaveTodos(key string, list Todos) error {
	if key == "" {
		return errors.New("Cannot save todo list with an empty key")
	}
	s.Lock()
	defer s.Unlock()
	s.repo[key] = list
	return nil
}

func CreateTodoCookie(w http.ResponseWriter) string {
	key := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	cookie := http.Cookie{
		Name:    CookieName,
		Value:   key,
		Expires: time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	return key
}
