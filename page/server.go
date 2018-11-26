package page

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rur/todowithoutmvc"
	"github.com/rur/treetop"
)

var CookieName = "todowithoutmvc-session"

type TodoHandler func(todowithoutmvc.Todos, treetop.Response, *http.Request) interface{}

type Server interface {
	Bind(TodoHandler) treetop.HandlerFunc
	LoadTodos(*http.Request) (todowithoutmvc.Todos, string)
	SaveTodos(string, todowithoutmvc.Todos) error
}

func NewServer(repo todowithoutmvc.Repository) Server {
	return &server{
		repo: repo,
	}
}

type server struct {
	sync.RWMutex
	repo todowithoutmvc.Repository
}

func (s *server) Bind(f TodoHandler) treetop.HandlerFunc {
	return func(rsp treetop.Response, req *http.Request) interface{} {
		// Here the Treetop response ID is being used to permit resources to be shared
		// between data handlers, within the scope of a request.
		todos, _ := s.LoadTodos(req)
		return f(todos, rsp, req)
	}
}

func (s *server) LoadTodos(req *http.Request) (todowithoutmvc.Todos, string) {
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

func (s *server) SaveTodos(key string, list todowithoutmvc.Todos) error {
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
