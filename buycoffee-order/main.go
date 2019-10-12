package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anubhavmishra/buycoffee/buycoffee-order/order"
	"github.com/braintree/manners"
	_ "github.com/lib/pq"
)

const version = "0.0.1"

func main() {

	var httpBindAddr = "0.0.0.0"
	var httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8082"
	}
	httpAddr := fmt.Sprintf("%s:%s", httpBindAddr, httpPort)

	// Setup postgres database access and configuration
	connStr := os.Getenv("POSTGRES_CONNECTION_STRING")
	if connStr == "" {
		log.Fatalln("postgres connection string \"POSTGRES_CONNECTION_STRING\" is required.")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening database: %s\n", err)
	}

	defer db.Close()

	// Check database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to database: %s\n", err)
	}

	log.Println("successfully connected to postgres database!")
	log.Println("Starting buycoffee order service.....")

	mux := http.NewServeMux()
	mux.Handle("/order", order.OrderHandler(db))
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
