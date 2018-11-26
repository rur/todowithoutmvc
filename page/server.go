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

type Server interface {
	Bind(ResourcesHandler) treetop.HandlerFunc
	LoadTodos(*http.Request) (todowithoutmvc.Todos, string)
	SaveTodos(string, todowithoutmvc.Todos) error
}

func NewServer(repo todowithoutmvc.Repository) Server {
	return &server{
		responses: make(map[uint32]*Resources),
		repo:      repo,
	}
}

type server struct {
	sync.RWMutex
	responses map[uint32]*Resources
	repo      todowithoutmvc.Repository
}

func (s *server) Bind(f ResourcesHandler) treetop.HandlerFunc {
	return func(rsp treetop.Response, req *http.Request) interface{} {
		// Here the Treetop response ID is being used to permit resources to be shared
		// between data handlers, within the scope of a request.
		respId := rsp.ResponseID()
		rsc := s.getResources(respId)

		if rsc == nil {
			todos, key := s.LoadTodos(req)
			if key == "" {
				// no todos key has been recorded, set a new cookie
				key = CreateTodoCookie(rsp)
			}

			rsc = &Resources{todos}

			s.setResources(respId, rsc)
			go func() {
				<-rsp.Context().Done()
				// assume that the request lifecycle is finished, just free up resources
				s.deleteResources(respId)
			}()
		}
		return f(*rsc, rsp, req)
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

// attempt to load Resources from the server cache
func (s *server) getResources(respId uint32) *Resources {
	s.RLock()
	defer s.RUnlock()
	if rsc, ok := s.responses[respId]; ok {
		return rsc
	} else {
		return nil
	}
}

func (s *server) setResources(respId uint32, rsc *Resources) {
	s.Lock()
	defer s.Unlock()
	s.responses[respId] = rsc
}

// remove Resources from the cache for a given treetop response ID, delete is idempotent
func (s *server) deleteResources(respId uint32) {
	s.Lock()
	defer s.Unlock()
	delete(s.responses, respId)
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
