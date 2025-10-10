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

// PublicComment is the public-facing comment structure without sensitive data
type PublicComment struct {
	ID        int    `json:"id"`
	ParentID  *int   `json:"parent_id,omitempty"`
	Nickname  string `json:"nickname"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// ToPublic converts a Comment to PublicComment
func (c *Comment) ToPublic() PublicComment {
	return PublicComment{
		ID:        c.ID,
		ParentID:  c.ParentID,
		Nickname:  c.Nickname,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}
