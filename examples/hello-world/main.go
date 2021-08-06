package main

import (
	"log"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: &Server{},
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

type Server struct{}

func (s *Server) ServeHTTP(http.ResponseWriter, *http.Request) {

}
