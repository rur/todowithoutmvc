package todowithoutmvc

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rur/treetop"
)

var CookieName = "todowithoutmvc-session"

type TodoHandler func(Todos, treetop.Response, *http.Request) interface{}

type Server interface {
	Bind(TodoHandler) treetop.HandlerFunc
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

func (s *server) Bind(f TodoHandler) treetop.HandlerFunc {
	return func(rsp treetop.Response, req *http.Request) interface{} {
		// load user todo list from repo based upon request cookies
		// pass to handler.
		//
		// Note that these handlers have no way to making changes
		// to the todo list.
		todos, _ := s.LoadTodos(req)
		return f(todos, rsp, req)
	}
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
