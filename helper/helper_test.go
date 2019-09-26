package helper

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/fshahy/sampleblog/article"

	_ "github.com/lib/pq"
)

var dataSourceTest = fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=disable",
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PASSWORD"),
	os.Getenv("POSTGRES_TEST_DB"),
)

func setup() {
	db, err := sql.Open("postgres", dataSourceTest)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	sqlText := `
		DROP TABLE public.articles;
		DROP SEQUENCE public.articles_id_seq;

		CREATE SEQUENCE public.articles_id_seq;

		CREATE TABLE public.articles (
			id INTEGER DEFAULT NEXTVAL('articles_id_seq'::regclass) PRIMARY KEY,
			title CHARACTER VARYING(255),
			"content" TEXT,
			author CHARACTER VARYING(255)
		);

		INSERT INTO articles(title, "content", author) VALUES('Post 1', 'Some interesting content goes here.', 'Farid');
		INSERT INTO articles(title, "content", author) VALUES('Post 2', 'Again some interesting content goes here.', 'Farid');
		INSERT INTO articles(title, "content", author) VALUES('Post 3', 'Yet more interesting content goes here.', 'Farid');`

	_, err = db.Exec(sqlText)
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	db, err := sql.Open("postgres", dataSourceTest)
	if err != nil {
		log.Panic(err)
	}

	sqlDelete := "DELETE FROM articles;"

	db.Exec(sqlDelete)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewResponseSlice(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()
	artcl, err := article.New("New Article", "Sample content for testing purpose.", "Developer")
	if err != nil {
		t.Fatal(err)
	}

	data := []article.Article{artcl}

	resp, err := NewResponse(200, "Success", data)
	if err != nil {
		t.Fatal(err)
	}
	if *resp.Status != 200 {
		t.Fatal("Error in creating new response status.")
	}

	if *resp.Message != "Success" {
		t.Fatal("Error in creating new response message.")
	}

	t.Log(reflect.TypeOf(data).Kind())

	if reflect.TypeOf(data).Kind() != reflect.Slice {
		t.Fatal("Error in creating new response data.")
	}

	tearDown()
}

func TestNewResponseInt64(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()
	// artcl := article.New("New Article", "Sample content for testing purpose.", "Developer")

	data := int64(1) //[]article.Article{artcl}

	resp, err := NewResponse(200, "Success", data)
	if err != nil {
		t.Fatal(err)
	}
	if *resp.Status != 200 {
		t.Fatal("Error in creating new response status.")
	}

	if *resp.Message != "Success" {
		t.Fatal("Error in creating new response message.")
	}

	t.Log(reflect.TypeOf(data).Kind())

	if reflect.TypeOf(data).Kind() != reflect.Int64 {
		t.Fatal("Error in creating new response data.")
	}

	tearDown()
}
