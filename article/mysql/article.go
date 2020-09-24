package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

type Article struct {
	UserID      int64
	ArticleID   int64
	ArticleName string
	Author      string
	Content     string
}

const (
	mysqlArticleCreateDB = iota
	mysqlArticleCreateTable
	mysqlArticleInsert
	mysqlArticleDeleteByID
	mysqlArticleUpdateByID
	mysqlArticleQueryByID
)

var (
	errInvalidInsert = errors.New("insert article:insert affected 0 rows")

	articleSQLstring = []string{
		`CREATE DATABASE IF NOT EXISTS %s`,
		`CREATE TABLE IF NOT EXISTS %s (
			articleId    INT(64)       NOT NULL AUTO_INCREMENT,
			userId       INT(64)       NOT NULL,
			articleName  VARCHAR(128)  NOT NULL,
			author       VARCHAR(128)  NOT NULL,
			content      TEXT          NOT NULL,
			PRIMARY KEY (articleId)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s (userId,articleName,author,content) VALUES (?,?,?,?)`,
		`DELETE FROM %s WHERE articleId = ? LIMIT 1`,
		// `UPDATE %s SET articleName = ? WHERE id = ? LIMIT 1`
		// `UPDATE %s SET author = ? WHERE id = ? LIMIT 1`,
		`UPDATE %s SET content = ? WHERE articleId = ? LIMIT 1`,
		`SELECT articleId, articleName, author, content FROM %s WHERE articleId = ? LIMIT 1`,
	}
)

//createDB
func CreateDB(db *sql.DB, DBName string) error {
	sql := fmt.Sprintf(articleSQLstring[mysqlArticleCreateDB], DBName)
	_, err := db.Exec(sql)
	return err
}

//createTable
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(articleSQLstring[mysqlArticleCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

//createArticle
func CreateArticle(db *sql.DB, tableName string, userId int, articleName string, author string, content string) (int, error) {
	sql := fmt.Sprintf(articleSQLstring[mysqlArticleInsert], tableName)
	result, err := db.Exec(sql, userId, articleName, author, content)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	articleId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(articleId), nil
}

//deleteArticleByID
func DeleteArticleByID(db *sql.DB, tableName string, id int) error {
	sql := fmt.Sprintf(articleSQLstring[mysqlArticleDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	return err
}

//updateArticleByID
func UpdateArticleByID(db *sql.DB, tableName string, content string, id int) error {
	sql := fmt.Sprintf(articleSQLstring[mysqlArticleUpdateByID], tableName)
	_, err := db.Exec(sql, content, id)
	return err
}

//queryArticleByID
func QueryArticleByID(db *sql.DB, tableName string, id int) (*Article, error) {
	var art Article

	sql := fmt.Sprintf(articleSQLstring[mysqlArticleQueryByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Content); err != nil {
			return nil, err
		}
	}

	return &art, nil
}
