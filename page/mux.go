package page

import (
	"net/http"
)

// Interface for router configuration
type Mux interface {
	Handle(pattern string, handler http.Handler)
}
