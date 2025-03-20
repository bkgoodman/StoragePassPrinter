package main

import (
	"fmt"
	"net/http"
    "sync"
    "context"
    "time"
)

// RequestData struct to hold the extracted parameters
type RequestData struct {
	Number string
	Name   string
}

// StartHTTPServer starts an HTTP daemon on the given address and port,
// and sends extracted request parameters to the provided channel.

func StartHTTPServer(addr string, requestChan chan<- RequestData) (func() error, error) {
    var wg sync.WaitGroup
    srv := &http.Server{Addr: addr}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        wg.Add(1)
        defer wg.Done()

        // Parse the request parameters
        err := r.ParseForm()
        if err != nil {
            fmt.Fprintf(w, "Error parsing form: %v\n", err)
            return
        }

        // Extract the 'number' and 'name' parameters
        number := r.Form.Get("number")
        name := r.Form.Get("name")

        // Create a RequestData struct
        data := RequestData{
            Number: number,
            Name:   name,
        }

        // Send the data to the channel
        select {
        case requestChan <- data:
            // Successfully sent
        case <-r.Context().Done():
            // Client disconnected or server is shutting down
            fmt.Println("Client disconnected or server shutting down during request processing")
            return
        }

        // Respond to the client (optional)
        fmt.Fprintf(w, "Request received and parameters sent to channel\n")
    })

    fmt.Printf("Starting HTTP server on %s\n", addr)
    go func() {
        err := srv.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            fmt.Printf("Error starting HTTP server: %v\n", err)
        } else {
            fmt.Println("HTTP server stopped")
        }
    }()

    // Return a function to shut down the server gracefully
    shutdown := func() error {
        fmt.Println("Shutting down HTTP server...")
        // Create a context with a timeout for shutdown
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Adjust timeout as needed
        defer cancel()

        // Initiate graceful shutdown
        if err := srv.Shutdown(ctx); err != nil {
            fmt.Printf("HTTP server shutdown error: %v\n", err)
            return err
        }

        // Wait for all pending requests to complete
        wg.Wait()
        fmt.Println("HTTP server gracefully stopped")
        return nil
    }

    return shutdown, nil
}
// Example usage (you would typically run this in a separate main function or test)
func http_init() (func() error) {
 	addr := "localhost:8080"
 	requestChannel := make(chan RequestData)

    shutdownHttp, err := StartHTTPServer(addr, requestChannel)
    if err != nil {
        fmt.Println("Server error:", err)
    }
    return shutdownHttp
 }
