package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// logging middleware to log all reqs.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		log.Printf("[%s] %s %s %d %s",
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
			sw.status,
			time.Since(start),
		)
	})
}

func main() {
    initDB()
    loadTemplates()

    mux := http.NewServeMux()

    // routes
    mux.HandleFunc("/", rootHandler)
    mux.HandleFunc("/post", postHandler)
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    loggedMux := loggingMiddleware(mux)

    fmt.Println("Guestbook is now running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
