package server

import (
	"fmt"
	"net/http"

	"github.com/jk0l3ti/cache-sse/cache"
)

type Server struct {
	Cache  cache.Cache
	Server http.Server
}

func readAndPushSse(cache cache.Cache, key string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}
		ch := make(chan any, 1)
		go func() {
			defer close(ch)
			cache.Stream(r.Context(), key, ch)
		}()
		for msg := range ch {
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush() // Immediately send data to the client
		}
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func NewServer(cacheType cache.CacheType) (*Server, error) {
	cache, err := cache.NewCache(cacheType)
	if err != nil {
		return nil, err
	}
	http.HandleFunc("/sse", readAndPushSse(cache, "name"))
	http.HandleFunc("/hello", handleDefault)
	return &Server{
		Cache: cache,
		Server: http.Server{
			Addr: ":8080",
		},
	}, nil
}

func (s *Server) Start() error {
	return s.Server.ListenAndServe()
}
