package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anubhavmishra/buycoffee/buycoffee-inventory/inventory"
	"github.com/braintree/manners"
)

const version = "0.0.1"

func main() {

	var httpBindAddr = "0.0.0.0"
	var httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8081"
	}
	httpAddr := fmt.Sprintf("%s:%s", httpBindAddr, httpPort)
	log.Println("Starting buycoffee inventory service.....")

	mux := http.NewServeMux()
	mux.HandleFunc("/restock", inventory.RestockHandler)
	mux.HandleFunc("/inventory", inventory.InventoryHandler)
	mux.HandleFunc("/health", HealthHandler)

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}

// LoggingHandler helps log web requests to the service
func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		format := "%s - - [%s] \"%s %s %s\" %s %s\n"
		fmt.Printf(format, r.RemoteAddr, time.Now().Format(time.RFC1123),
			r.Method, r.URL.Path, r.Proto, r.Header.Get("X-Service-ID"), r.UserAgent())
		h.ServeHTTP(w, r)
	})
}
