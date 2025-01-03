package repositories

import (
	"database/sql"
	"fmt"
	"go-crud-sandbox/models"
)

const articleNumPerPage = 5

func InsertArticle(db *sql.DB, article models.Article) (models.Article, error) {
	const sqlStr = `insert into articles (title, contents, username, nice, created_at) values(?, ?, ?, ?, ?);`

	result, err := db.Exec(sqlStr, article.Title, article.Contents, article.UserName, article.NiceNum, article.CreatedAt)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to insert article: %w", err)
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		return models.Article{}, err
	}

	newArticle := models.Article{
		ID:        int(articleID),
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum,
		CreatedAt: article.CreatedAt,
	}

	return newArticle, nil
}

func SelectArticleList(db *sql.DB, page int) ([]models.Article, error) {
	const sqlStr = `select article_id, title, contents, username, nice from articles limit ? offset ?;`
	rows, err := db.Query(sqlStr, articleNumPerPage, ((page - 1) * articleNumPerPage))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articleArray []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum)
		if err != nil {
			return nil, err
		}
		articleArray = append(articleArray, article)
	}

	return articleArray, nil
}

func SelectArticleDetail(db *sql.DB, articleID int) (models.Article, error) {
	const sqlStr = `select * from articles where article_id = ?;`
	row := db.QueryRow(sqlStr, articleID)
	if err := row.Err(); err != nil {
		fmt.Println(err)
		return models.Article{}, err
	}

	var article models.Article
	var createdTime sql.NullTime
	err := row.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &article.CreatedAt)
	if err != nil {
		return models.Article{}, err
	}

	if createdTime.Valid {
		article.CreatedAt = createdTime.Time
	}

	return article, nil
}

func UpdateNiceNum(db *sql.DB, articleID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	const sqlGetNice = `select nice from articles where article_id = ?;`
	row := tx.QueryRow(sqlGetNice, articleID)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}

	var niceNum int
	err = row.Scan(&niceNum)
	if err != nil {
		tx.Rollback()
		return err
	}

	const sqlUpdateNice = `update articles set nice = ? where article_id = ?`
	_, err = tx.Exec(sqlUpdateNice, niceNum+1, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
