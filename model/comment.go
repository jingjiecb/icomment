package model

type Comment struct {
	ID         int    `json:"id"`
	ArticleURL string `json:"article_url"`
	ParentID   *int   `json:"parent_id,omitempty"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Content    string `json:"content"`
	Status     string `json:"status"` // "pending" or "approved"
	CreatedAt  string `json:"created_at"`
}
