package dao

import (
	"testing"

	"claws.top/icomment/model"
)

func TestUpsertComment(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()
	RunSqlFile(t, db, "../sql/create_table_comment.sql")
	dao := NewCommentDao(db)

	tests := []struct {
		name         string
		inputComment *model.Comment
		expectErr    bool
		expectOutput *model.Comment
	}{
		{
			name: "insert comment without id",
			inputComment: &model.Comment{
				Author:  "test@example.com",
				Content: "mock content",
			},
			expectErr: false,
			expectOutput: &model.Comment{
				ID:      1,
				Author:  "test@example.com",
				Content: "mock content",
			},
		},
		{
			name: "update comment with existing id",
			inputComment: &model.Comment{
				ID:      1,
				Author:  "test-new@example.com",
				Content: "updated content",
			},
			expectErr: false,
			expectOutput: &model.Comment{
				ID:      1,
				Author:  "test-new@example.com",
				Content: "updated content",
			},
		},
		{
			name: "Return error when given not existing id",
			inputComment: &model.Comment{
				ID:      999,
				Author:  "test-new@example.com",
				Content: "updated content",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedComment, err := dao.UpsertComment(tt.inputComment)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expect error: %v, got error: %v", tt.expectErr, err)
			}
			if !tt.expectErr {
				if updatedComment.CreatedAt == "" {
					t.Fatalf("expect created time should not be empty, but got empty one")
				}
				if tt.expectOutput.ID != updatedComment.ID || tt.expectOutput.Author != updatedComment.Author || tt.expectOutput.Content != updatedComment.Content {
					t.Fatalf("expect id, author, content should be the same. expected: %v, actual: %v", tt.expectOutput, updatedComment)
				}

			}
		})
	}
}
