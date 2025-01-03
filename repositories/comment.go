package repositories

import (
	"database/sql"
	"go-crud-sandbox/models"
)

func InsertComment(db *sql.DB, comment models.Comment) (models.Comment, error) {
	const sqlStr = `insert into comments (article_id, message, created_at) values (?, ?, ?);`
	result, err := db.Exec(sqlStr, comment.ArticleID, comment.Message, comment.CreatedAt)
	if err != nil {
		return models.Comment{}, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, err
	}
	newComment := models.Comment{
		ID:        int(commentID),
		ArticleID: comment.ArticleID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}

	return newComment, nil
}

func SelectCommentList(db *sql.DB, articleID int) ([]models.Comment, error) {
	const sqlStr = `select * from comments where article_id = ?;`
	rows, err := db.Query(sqlStr, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentArray []models.Comment
	for rows.Next() {
		var comment models.Comment
		var createdTime sql.NullTime
		err := rows.Scan(&comment.ID, &comment.ArticleID, &comment.Message, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		if createdTime.Valid {
			comment.CreatedAt = createdTime.Time
		}
		commentArray = append(commentArray, comment)
	}

	return commentArray, nil
}
