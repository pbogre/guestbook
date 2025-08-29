package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "golang.org/x/time/rate"
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

// ratelimit middleware
func rateLimitMiddleware(next http.Handler) http.Handler {
    rateLimit := 1.0 / float64(GuestbookConfig.EntriesPerPage)
    burstLimit := GuestbookConfig.GlobalBurstLimit
    limiter := rate.NewLimiter(rate.Limit(rateLimit), burstLimit)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "The server is handling too many requests, please wait a few minutes", http.StatusTooManyRequests)
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    initDB("./data/guestbook.db")
    loadConfig("./config/guestbook.yml")
    loadTemplates()

    mux := http.NewServeMux()

    // routes
    mux.HandleFunc("/", rootHandler)
    mux.Handle("/post", rateLimitMiddleware(http.HandlerFunc(postHandler)))
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    loggedMux := loggingMiddleware(mux)

    fmt.Println("Guestbook is now running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
