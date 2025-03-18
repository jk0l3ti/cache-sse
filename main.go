package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jk0l3ti/cache-sse/cache"
	"github.com/jk0l3ti/cache-sse/server"
)

func main() {
	cacheType := cache.Redis
	server, err := server.NewServer(cacheType)
	if err != nil {
		log.Fatal("failed to start server", err.Error())
		return
	}
	fmt.Printf("got %v cache connection, starting SSE server\n", cacheType)
	err = server.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server stopped with err: ", err.Error())
	}
}
