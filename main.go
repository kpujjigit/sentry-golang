// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Sentry with tracesSampleRate and EnableProfiling
	clientOptions := sentry.ClientOptions{
		Dsn:                os.Getenv("SENTRY_DSN"),     // Set DSN from environment variable
		Release:            os.Getenv("SENTRY_RELEASE"), // Set release from environment variable
		EnableTracing:      true,
		TracesSampleRate:   1.0, // Capture 100% of transactions
		ProfilesSampleRate: 1.0, // Enable profiling
	}

	if clientOptions.Release == "" {
		log.Fatalf("sentry.Init: Release identifier not set")
	}

	err = sentry.Init(clientOptions)
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create a Sentry HTTP handler
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Step 3: Create a simple HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Start a root transaction
		tx := sentry.StartTransaction(r.Context(), "GET /")
		defer tx.Finish()

		// Create a custom span
		span := tx.StartChild("custom.operation")
		time.Sleep(100 * time.Millisecond) // Simulate some work
		span.Finish()

		// Add a breadcrumb
		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Message: "User visited the homepage",
			Level:   sentry.LevelInfo,
		})

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

		// Custom instrumentation
		childSpan1 := span.StartChild("operation1")
		time.Sleep(1 * time.Second) // Simulate some work
		childSpan1.Finish()

		childSpan2 := span.StartChild("operation2")
		time.Sleep(1 * time.Second) // Simulate some more work
		childSpan2.Finish()

		time.Sleep(2 * time.Second) // Simulate a delay
		w.Write([]byte("Performance endpoint"))
	})

	http.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		// Collect user feedback
		sentry.CaptureMessage("User feedback collected")
		w.Write([]byte("Thank you for your feedback!"))
	})

	http.HandleFunc("/custom-tags", func(w http.ResponseWriter, r *http.Request) {
		// Add custom tags
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("custom-tag", "example")
		})
		sentry.CaptureMessage("Custom tags added")
		w.Write([]byte("Custom tags added"))
	})

	http.HandleFunc("/context", func(w http.ResponseWriter, r *http.Request) {
		// Add additional context information
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetContext("example", map[string]interface{}{
				"key": "value",
			})
		})
		sentry.CaptureMessage("Context information added")
		w.Write([]byte("Context information added"))
	})

	// Wrap the mux with the Sentry handler
	http.Handle("/", sentryHandler.Handle(mux))

	// Step 5: Run the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
