package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"claws.top/icomment/dao"
)

//go:embed static/*
var staticFS embed.FS

//go:embed templates/*
var templateFS embed.FS

func main() {
	cfg := LoadConfig()

	db, err := InitDB(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	commentDao := dao.NewCommentDao(db)

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"eq":  func(a, b string) bool { return a == b },
	}

	tpl, err := template.New("admin.html").Funcs(funcMap).ParseFS(templateFS, "templates/admin.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}

	// Initialize Bark notifier if configured
	barkNotifier := NewBarkNotifier(cfg.BarkDeviceKey)
	if barkNotifier != nil {
		fmt.Printf("   Bark:     Enabled\n")
	}

	server := NewServer(commentDao, tpl, barkNotifier)

	// Public API server
	publicMux := http.NewServeMux()

	// RESTful API endpoints
	// GET /api/comments?article_url=xxx - List comments for an article
	// POST /api/comments - Create a new comment
	publicMux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.GetComments(w, r)
		case http.MethodPost:
			server.CreateComment(w, r)
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Static files
	publicMux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		content, err := staticFS.ReadFile(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
			// Cache for 1 hour (3600 seconds)
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		w.Write(content)
	})

	// Admin server - RESTful design
	adminMux := http.NewServeMux()

	// GET /comments - List/filter comments (with pagination)
	// DELETE /comments/:id - Delete a comment
	// PATCH /comments/:id/approve - Approve a comment
	adminMux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			server.AdminListComments(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	adminMux.HandleFunc("/comments/", func(w http.ResponseWriter, r *http.Request) {
		server.AdminCommentActions(w, r)
	})

	// GET / - Admin UI page
	adminMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			server.AdminPage(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	publicAddr := ":" + cfg.Port
	adminAddr := ":" + cfg.AdminPort

	fmt.Printf("🚀 iComment server started\n")
	fmt.Printf("   Database: %s\n", cfg.DBPath)
	fmt.Printf("   Public:   http://0.0.0.0%s\n", publicAddr)
	fmt.Printf("   Admin:    http://0.0.0.0:%s (WARNING: Exposed on all interfaces)\n", cfg.AdminPort)

	// Start admin server in goroutine
	go func() {
		if err := http.ListenAndServe(adminAddr, adminMux); err != nil {
			log.Fatal("Admin server failed:", err)
		}
	}()

	// Start public server
	if err := http.ListenAndServe(publicAddr, publicMux); err != nil {
		log.Fatal("Public server failed:", err)
	}
}
