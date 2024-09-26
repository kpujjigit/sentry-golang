// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func main() {

	// Step 4: Initialize Sentry with tracesSampleRate and EnableProfiling
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                "https://ad297c80676444d7bf21a3919c2b6d5a@o4504052292517888.ingest.us.sentry.io/4504300727566336", // Replace with your Sentry DSN
		EnableTracing:      true,
		TracesSampleRate:   1.0, // Capture 100% of transactions
		ProfilesSampleRate: 1.0, // Enable profiling
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create a Sentry HTTP handler
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Step 3: Create a simple HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		// Simulate an error
		err := fmt.Errorf("something went wrong")
		sentry.CaptureException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})

	http.HandleFunc("/performance", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a performance issue
		span := sentry.StartSpan(r.Context(), "performance")
		defer span.Finish()

		time.Sleep(2 * time.Second) // Simulate a delay
		w.Write([]byte("Performance endpoint"))
	})

	// Wrap the mux with the Sentry handler
	http.Handle("/", sentryHandler.Handle(mux))

	// Step 5: Run the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
