/*
Package article implements functionalities for blog article.
An article has the following fields:
	ID
	Title
	Content
	Author
*/
package article

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// a handle to database connection pool.
var db *sql.DB

// ConnectToDB opens the database connection and sets db variable.
func ConnectToDB(dataSource string) {
	var err error
	db, err = sql.Open("postgres", dataSource)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
}

// CloseDB closes database connection.
func CloseDB() {
	db.Close()
}

// Article is a type representing blog articles.
type Article struct {
	ID      *int64  `json:"id,omitempty"`
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Author  *string `json:"author,omitempty"`
}

// New creates a new blog article ans saves it to database.
func New(title string, content string, author string) (Article, error) {
	var articleID int64
	sqlCreate := "INSERT INTO articles(title, content, author) VALUES($1, $2, $3) RETURNING id"

	err := db.QueryRow(sqlCreate, title, content, author).Scan(&articleID)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:      &articleID,
		Title:   &title,
		Content: &content,
		Author:  &author,
	}, nil
}

// FindByID searches database for an article with specified id.
func FindByID(id int64) (Article, error) {
	sqlFind := "SELECT title, content, author FROM articles WHERE id = $1"
	rows, err := db.Query(sqlFind, id)
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}

	if rows.Next() {
		var rowTitle, rowContent, rowAuthor string
		err := rows.Scan(&rowTitle, &rowContent, &rowAuthor)
		if err != nil {
			log.Println(err)
		}

		return Article{
			ID:      &id,
			Title:   &rowTitle,
			Content: &rowContent,
			Author:  &rowAuthor,
		}, nil
	}

	errMessage := fmt.Sprintf("no article with id=%d exists", id)
	err = errors.New(errMessage)
	return Article{}, err
}

// GetAll returns all exisiting articles.
// returns an slice of Articles
func GetAll() ([]Article, error) {
	var articles []Article

	sqlFindAll := `SELECT id, title, content, author FROM articles`
	rows, err := db.Query(sqlFindAll)
	defer rows.Close()
	if err != nil {
		return []Article{}, err
	}

	for rows.Next() {
		artcl := new(Article)
		err := rows.Scan(&artcl.ID, &artcl.Title, &artcl.Content, &artcl.Author)
		if err != nil {
			log.Println(err)
		}

		articles = append(articles, *artcl)
	}

	return articles, nil
}
