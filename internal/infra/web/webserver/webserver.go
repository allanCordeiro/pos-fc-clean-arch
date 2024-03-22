package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Route struct {
	HTTPMethod  string
	HandlerFunc http.HandlerFunc
	Path        string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]Route
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]Route),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(httpMethod string, path string, handler http.HandlerFunc) {
	s.Handlers[fmt.Sprintf("%s/%s", httpMethod, path)] = Route{HTTPMethod: httpMethod, HandlerFunc: handler, Path: path}
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handler := range s.Handlers {
		switch handler.HTTPMethod {
		case "POST":
			s.Router.Post(handler.Path, handler.HandlerFunc)
		case "GET":
			s.Router.Get(handler.Path, handler.HandlerFunc)
		default:
			panic("unrecognized method")
		}

		//s.Router.Handle(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
