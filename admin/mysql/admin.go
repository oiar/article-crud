package mysql

import (
	"database/sql"
	"errors"

	"github.com/yhyddr/article-crud/pkgs/hash"
)

const (
	mysqlUserCreateTable = iota
	mysqlUserInsert
	mysqlUserLogin
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed  = errors.New("invalid username or password")

	adminSQLString = []string{
		`CREATE TABLE IF NOT EXISTS article.admin (
			admin_id  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name     	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			password 	VARCHAR(512) NOT NULL DEFAULT ' ',
			PRIMARY KEY (admin_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO article.admin (name,password) VALUES (?,?)`,
		`SELECT admin_id,password FROM article.admin WHERE name = ? LOCK IN SHARE MODE`,
	}
)

//create admin table.
func CreateTable(db *sql.DB, name, password *string) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}

	Create(db, name, password)
	return nil
}

//create an user
func Create(db *sql.DB, name, password *string) error {
	hash, err := hash.Generate(password)
	if err != nil {
		return err
	}

	result, err := db.Exec(adminSQLString[mysqlUserInsert], name, hash)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//the user logins
func Login(db *sql.DB, name, password *string) (uint32, error) {
	var (
		id  uint32
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserLogin], name).Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}

	if !hash.Compare([]byte(pwd), password) {
		return 0, errLoginFailed
	}

	return id, nil
}
