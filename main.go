package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type IPInfo struct {
	IPv4      string
	IPv6      string
	UserAgent string
}

type DebugInfo struct {
	Headers    map[string][]string
	Host       string
	RemoteAddr string
}

var templates map[string]*template.Template

func initTemplates() error {
	templates = make(map[string]*template.Template)

	templateFiles := map[string]string{
		"index":      filepath.Join("templates", "index.html.tmpl"),
		"plain":      filepath.Join("templates", "plain.txt.tmpl"),
		"api":        filepath.Join("templates", "api.json.tmpl"),
		"debugHtml":  filepath.Join("templates", "debug.html.tmpl"),
		"debugPlain": filepath.Join("templates", "debug.txt.tmpl"),
		"debugJson":  filepath.Join("templates", "debug.json.tmpl"),
	}

	for name, file := range templateFiles {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", name, err)
		}
		templates[name] = tmpl
	}

	return nil
}

func main() {
	if err := initTemplates(); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if isPlainTextRequest(r) {
			http.Redirect(w, r, "/plain", http.StatusSeeOther)
			return
		}

		ipInfo := getIPInfo(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates["index"].Execute(w, ipInfo)
	})

	http.HandleFunc("/plain", func(w http.ResponseWriter, r *http.Request) {
		ipInfo := getIPInfo(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		templates["plain"].Execute(w, ipInfo)
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		ipInfo := getIPInfo(r)
		w.Header().Set("Content-Type", "application/json")
		templates["api"].Execute(w, ipInfo)
	})

	http.HandleFunc("/ipv4", func(w http.ResponseWriter, r *http.Request) {
		ipv4 := getIP(r)
		if !isIPv4(ipv4) {
			ipv4 = "Not available"
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, ipv4)
	})

	http.HandleFunc("/ipv6", func(w http.ResponseWriter, r *http.Request) {
		ipv6 := getIP(r)
		if !isIPv6(ipv6) {
			ipv6 = "Not available"
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, ipv6)
	})

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		debugInfo := getDebugInfo(r)

		accept := r.Header.Get("Accept")
		isPlain := isPlainTextRequest(r)

		switch {
		case strings.Contains(accept, "application/json"):
			w.Header().Set("Content-Type", "application/json")
			templates["debugJson"].Execute(w, debugInfo)
		case isPlain || strings.Contains(accept, "text/plain"):
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			templates["debugPlain"].Execute(w, debugInfo)
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			templates["debugHtml"].Execute(w, debugInfo)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getIPInfo(r *http.Request) IPInfo {
	directIP := getIP(r)

	var ipv4, ipv6 string

	if isIPv4(directIP) {
		ipv4 = directIP
		ipv6 = "Not available"
	} else if isIPv6(directIP) {
		ipv6 = directIP
		ipv4 = "Not available"
	} else {
		ipv4 = "Not available"
		ipv6 = "Not available"
	}

	return IPInfo{
		IPv4:      ipv4,
		IPv6:      ipv6,
		UserAgent: r.UserAgent(),
	}
}

func getDebugInfo(r *http.Request) DebugInfo {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	headers := make(map[string][]string)
	for name, values := range r.Header {
		headers[name] = values
	}

	return DebugInfo{
		Headers:    headers,
		Host:       hostname,
		RemoteAddr: r.RemoteAddr,
	}
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func isIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To4() != nil
}

func isIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To4() == nil
}

func isPlainTextRequest(r *http.Request) bool {
	userAgent := strings.ToLower(r.UserAgent())
	if strings.Contains(userAgent, "curl") ||
		strings.Contains(userAgent, "wget") ||
		strings.Contains(userAgent, "httpie") {
		return true
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "text/plain") &&
		!strings.Contains(accept, "text/html") {
		return true
	}

	return false
}
