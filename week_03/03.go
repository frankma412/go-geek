package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewHttpServer(out chan struct{}) http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	//mock: http server stop
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		out <- struct{}{}
	})

	return http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
}

func main() {
	//init
	g, ctx := errgroup.WithContext(context.Background())
	out := make(chan struct{})
	httpServer := NewHttpServer(out)

	//run httpserver
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	//signal server
	g.Go(func() error {
		//mock: stop signal
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGUSR1)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case s := <-q:
				return errors.Errorf("os signal: %v", s)
			}
		}
	})

	//server out
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				log.Println("errgroup exit")
			case <-out:
				log.Println("server outing")
			}

			toCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			log.Println("shutdown server")
			return httpServer.Shutdown(toCtx)
		}
	})

	//mock: error occurred
	g.Go(func() error {
		fmt.Println("mock error ..... ")
		time.Sleep(3 * time.Second)
		fmt.Println("error occurred")
		return errors.New("mock error occurred")
	})


	err := g.Wait()
	fmt.Printf("app exit: %v", err)
}
