package kagami

import (
	"fmt"
	"html"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Server is a struct holding the configuration of git-mirror
type Server struct {
}

func githubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// NewServer creates a new server instance that can answer to hooks
func NewServer(config *Config) (*Server, error) {
	mux := http.NewServeMux()

	// Install every handler
	for name, provider := range GetProviderInstances() {
		log.Infof("serving provider %s at /provider/%s", name, name)
		mux.Handle("/provider/"+name, provider)
	}

	log.Infof("starting http server on address %s", config.Server.Addr)
	err := http.ListenAndServe(config.Server.Addr, mux)
	if err != nil {
		return nil, err
	}
	return &Server{}, nil
}
