package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"claws.top/icomment/dao"
	"claws.top/icomment/model"
)

type Server struct {
	dao          *dao.CommentDao
	tpl          *template.Template
	barkNotifier *BarkNotifier
}

func NewServer(d *dao.CommentDao, tpl *template.Template, notifier *BarkNotifier) *Server {
	return &Server{
		dao:          d,
		tpl:          tpl,
		barkNotifier: notifier,
	}
}

// Helper functions
func (s *Server) jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Public API Handlers
// GET /api/comments?article_url=xxx
func (s *Server) GetComments(w http.ResponseWriter, r *http.Request) {
	articleURL := r.URL.Query().Get("article_url")
	if articleURL == "" {
		s.jsonError(w, "article_url parameter required", http.StatusBadRequest)
		return
	}

	comments, err := s.dao.GetCommentsByURL(articleURL)
	if err != nil {
		s.jsonError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert to public comments (remove sensitive data)
	publicComments := make([]model.PublicComment, len(comments))
	for i, c := range comments {
		publicComments[i] = c.ToPublic()
	}

	s.jsonResponse(w, publicComments, http.StatusOK)
}

// POST /api/comments
func (s *Server) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		s.jsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if comment.ArticleURL == "" || comment.Nickname == "" || comment.Content == "" {
		s.jsonError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate field lengths
	if len(comment.Nickname) > 50 || len(comment.Email) > 100 || len(comment.Content) > 2000 {
		s.jsonError(w, "Field length exceeds maximum allowed", http.StatusBadRequest)
		return
	}

	if err := s.dao.CreateComment(&comment); err != nil {
		if err == sql.ErrNoRows {
			s.jsonError(w, "Parent comment not found", http.StatusBadRequest)
		} else {
			s.jsonError(w, "Failed to create comment", http.StatusInternalServerError)
		}
		return
	}

	// Send Bark notification (non-blocking)
	if s.barkNotifier != nil {
		go func() {
			if err := s.barkNotifier.SendNewCommentNotification(&comment); err != nil {
				// Log error but don't fail the request
				println("Failed to send Bark notification:", err.Error())
			}
		}()
	}

	s.jsonResponse(w, map[string]string{"message": "Comment created, pending approval"}, http.StatusCreated)
}

// Admin Handlers
func (s *Server) AdminPage(w http.ResponseWriter, r *http.Request) {
	filter := s.parseFilter(r)
	comments, total, err := s.dao.GetCommentsWithFilter(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Comments":   comments,
		"Total":      total,
		"Page":       filter.Page,
		"TotalPages": (total + filter.PageSize - 1) / filter.PageSize,
		"Status":     filter.Status,
		"ArticleURL": filter.ArticleURL,
	}
	s.tpl.Execute(w, data)
}

// GET /comments - List/filter comments with pagination (JSON API)
func (s *Server) AdminListComments(w http.ResponseWriter, r *http.Request) {
	filter := s.parseFilter(r)
	comments, total, err := s.dao.GetCommentsWithFilter(filter)
	if err != nil {
		s.jsonError(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	s.jsonResponse(w, map[string]interface{}{
		"comments":    comments,
		"total":       total,
		"page":        filter.Page,
		"total_pages": (total + filter.PageSize - 1) / filter.PageSize,
		"page_size":   filter.PageSize,
	}, http.StatusOK)
}

func (s *Server) parseFilter(r *http.Request) dao.CommentFilter {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending"
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	return dao.CommentFilter{
		Status:     status,
		ArticleURL: r.URL.Query().Get("article_url"),
		Page:       page,
		PageSize:   10,
	}
}

// Handle /comments/:id and /comments/:id/approve
func (s *Server) AdminCommentActions(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// DELETE /comments/:id
	if len(parts) == 2 && r.Method == http.MethodDelete {
		if err := s.dao.DeleteComment(id); err != nil {
			s.jsonError(w, "Failed to delete comment", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// PATCH /comments/:id/approve
	if len(parts) == 3 && parts[2] == "approve" && r.Method == http.MethodPatch {
		if err := s.dao.ApproveComment(id); err != nil {
			s.jsonError(w, "Failed to approve comment", http.StatusInternalServerError)
			return
		}
		s.jsonResponse(w, map[string]string{"message": "Comment approved"}, http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
