package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
        "strconv"
        "net"

	// "github.com/dchest/captcha"
)

var templates *template.Template

func loadTemplates() {
    tmplFiles, _ := filepath.Glob("templates/*.html")

    // define custom template functions
    templates = template.New("").Funcs(template.FuncMap{
        "sub": func(a, b int) int { return a - b },
        "add": func(a, b int) int { return a + b },
    })

    // parse files
    templates = template.Must(templates.ParseFiles(tmplFiles...))
}

func stripRemoteAddr(remoteAddr string) (string, error) {
    // get only IP address (for unique messages) instead of IP:port
    remoteAddr, _, err := net.SplitHostPort(remoteAddr)
    return remoteAddr, err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    // CurrentPage
    currentPage, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        currentPage = 1
    }

    // TotalPages
    totalPages, err := getTotalPages() 
    if err != nil {
        http.Error(w, "Failed to get total pages", http.StatusInternalServerError)
        log.Print(err)
        return
    }

    if currentPage < 1 || currentPage > totalPages {
        http.Error(w, "Invalid page number", http.StatusBadRequest)
        return
    }

    // Messages
    messages, err := getMessages(currentPage - 1)
    if err != nil {
        http.Error(w, "Failed to load messages", http.StatusInternalServerError)
        log.Print(err)
        return
    }

    // CanWrite
    remoteAddr, err := stripRemoteAddr(r.RemoteAddr)
    if err != nil {
        http.Error(w, "Failed to save message", http.StatusInternalServerError)
        log.Print(err)
        return
    }
    canWrite, err := canRemoteAddrWrite(remoteAddr)

    data := map[string]any{
        "CurrentPage": currentPage,
        "TotalPages": totalPages,
        "Messages": messages,
        "CanWrite": canWrite,
    }
    templates.ExecuteTemplate(w, "layout.html", data)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // input sanitization
    name := strings.TrimSpace(r.FormValue("name"))
    content := strings.TrimSpace(r.FormValue("content"))

    name = html.EscapeString(name)
    content = html.EscapeString(content)

    if name == "" || content == "" {
        http.Error(w, "Name and content are required", http.StatusBadRequest)
        return
    }

    if len(name) > 22 || len(content) > 250 {
        http.Error(w, "Name or content exceed maximum length", http.StatusBadRequest)
        return
    }

    remoteAddr, err := stripRemoteAddr(r.RemoteAddr)
    if err != nil {
        http.Error(w, "Failed to save message", http.StatusInternalServerError)
        log.Print(err)
        return
    }

    if err := postMessage(name, content, remoteAddr); err != nil {
        // each user can only send one message (ip-based)
        if strings.Contains(err.Error(), "UNIQUE constraint failed") {
            http.Error(w, "You have already written a message", http.StatusTooManyRequests)
            log.Print(err)
            return
        }
        http.Error(w, "Failed to save message", http.StatusInternalServerError)
        log.Print(err)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
