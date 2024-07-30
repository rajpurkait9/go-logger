package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// CustomResponseWriter wraps http.ResponseWriter to capture the status code
type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	ErrorMsg   string
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware logs the incoming HTTP requests with color
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		crw := &CustomResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		methodColor := color.New(color.FgGreen).SprintFunc()
		urlColor := color.New(color.FgCyan).SprintFunc()
		next.ServeHTTP(crw, r)
		durationColor := color.New(color.FgYellow).SprintFunc()
		statusColor := color.New(color.FgRed).SprintFunc()
		if crw.StatusCode < 400 {
			statusColor = color.New(color.FgGreen).SprintFunc()
		}
		if crw.StatusCode >= 400 {
			log.Printf("%s %s %s %s in %s - Error: %s", methodColor(r.Method), urlColor(r.RequestURI), "Status-code", statusColor(crw.StatusCode), durationColor(time.Since(start)), crw.ErrorMsg)
		} else {
			log.Printf("%s %s %s %s in %s", methodColor(r.Method), urlColor(r.RequestURI), "Status-code", statusColor(crw.StatusCode), durationColor(time.Since(start)))
		}
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", HelloHandler)
	r.HandleFunc("/goodbye", GoodbyeHandler)

	// Apply the logging middleware
	r.Use(LoggingMiddleware)

	http.Handle("/", r)
	defineArr()
	defineStruct()
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world! from raj purkait the great man"))
}

func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Goodbye, world!"))
}

// Dummy functions to avoid compilation errors
func defineArr()    {}
func defineStruct() {}
