package dao

import (
	"database/sql"

	"claws.top/icomment/model"
)

type CommentDao struct {
	db *sql.DB
}

func NewCommentDao(db *sql.DB) *CommentDao {
	return &CommentDao{db: db}
}

func (dao *CommentDao) CreateComment(comment *model.Comment) error {
	// Validate parent_id exists if provided
	if comment.ParentID != nil {
		var exists bool
		err := dao.db.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?)", *comment.ParentID).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return sql.ErrNoRows // Will be handled as bad request
		}
	}

	_, err := dao.db.Exec(
		"INSERT INTO comments(article_url, parent_id, nickname, email, content) VALUES(?, ?, ?, ?, ?)",
		comment.ArticleURL, comment.ParentID, comment.Nickname, comment.Email, comment.Content,
	)
	return err
}

func (dao *CommentDao) GetCommentsByURL(url string) ([]model.Comment, error) {
	rows, err := dao.db.Query(
		"SELECT id, article_url, parent_id, nickname, email, content, status, created_at FROM comments WHERE article_url = ? AND status = 'approved' ORDER BY created_at ASC",
		url,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ArticleURL, &c.ParentID, &c.Nickname, &c.Email, &c.Content, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

type CommentFilter struct {
	Status     string // "all", "pending", "approved"
	ArticleURL string
	Email      string
	Page       int
	PageSize   int
}

func (dao *CommentDao) GetCommentsWithFilter(filter CommentFilter) ([]model.Comment, int, error) {
	query := "SELECT id, article_url, parent_id, nickname, email, content, status, created_at FROM comments WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM comments WHERE 1=1"
	args := []interface{}{}

	if filter.Status != "all" && filter.Status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, filter.Status)
	}

	if filter.ArticleURL != "" {
		query += " AND article_url LIKE ?"
		countQuery += " AND article_url LIKE ?"
		args = append(args, filter.ArticleURL+"%")
	}

	if filter.Email != "" {
		query += " AND email = ?"
		countQuery += " AND email = ?"
		args = append(args, filter.Email)
	}

	// Get total count
	var total int
	err := dao.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query += " ORDER BY created_at ASC LIMIT ? OFFSET ?"
	offset := (filter.Page - 1) * filter.PageSize
	queryArgs := append(args, filter.PageSize, offset)

	rows, err := dao.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ArticleURL, &c.ParentID, &c.Nickname, &c.Email, &c.Content, &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}
	return comments, total, nil
}

func (dao *CommentDao) ApproveComment(id int) error {
	_, err := dao.db.Exec("UPDATE comments SET status = 'approved' WHERE id = ?", id)
	return err
}

func (dao *CommentDao) DeleteComment(id int) error {
	_, err := dao.db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}
