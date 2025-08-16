package dao

import (
	"database/sql"
	"fmt"

	"claws.top/icomment/model"
)

type CommentDao struct {
	db *sql.DB
}

func NewCommentDao(db *sql.DB) *CommentDao {
	return &CommentDao{db: db}
}

func (dao *CommentDao) UpsertComment(comment *model.Comment) (*model.Comment, error) {
	if comment.ID == 0 {
		id, err := dao.insertComment(comment)
		if err != nil {
			return nil, err
		}

		comment, err := dao.getCommentById(id)
		if err != nil {
			return nil, err
		}

		return comment, nil
	}

	err := dao.updateComment(comment)
	if err != nil {
		return nil, err
	}

	comment, err = dao.getCommentById(comment.ID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (dao *CommentDao) insertComment(comment *model.Comment) (int, error) {
	res, err := dao.db.Exec("INSERT INTO comments(author, content) VALUES(?, ?)", comment.Author, comment.Content)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), err
}

func (dao *CommentDao) updateComment(comment *model.Comment) error {
	_, err := dao.db.Exec("UPDATE comments SET author=?, content=? WHERE id=?", comment.Author, comment.Content, comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func (dao *CommentDao) getCommentById(id int) (*model.Comment, error) {
	stmt, err := dao.db.Prepare("SELECT id, author, content, created_at FROM comments WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var c model.Comment
	err = stmt.QueryRow(id).Scan(&c.ID, &c.Author, &c.Content, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment with id %d not found", id)
		}
		return nil, err
	}

	return &c, nil
}
