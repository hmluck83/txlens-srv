package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const listenPort = ":7814"

type routeHttp struct {
	server *http.Server
}

func NewRouter() *routeHttp {
	return &routeHttp{
		server: &http.Server{
			Addr: listenPort,
		},
	}
}

func (rh *routeHttp) Run() {
	http.HandleFunc("/", summuryHandler)

	go func() {
		log.Printf("Server staring on %s\n", rh.server.Addr)
		if err := rh.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Cound not listen on %s: %v\n", rh.server.Addr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Context 리소스 해제

	if err := rh.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Panicln("Server exited gracefully")

}
