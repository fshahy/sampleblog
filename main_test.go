package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fshahy/sampleblog/article"
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

func TestGetArticlByID(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()

	req, err := http.NewRequest("GET", "/articles/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReq)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"Success","data":[{"id":1,"title":"Post 1","content":"Some interesting content goes here.","author":"Farid"}]}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

	tearDown()
}

func TestGetNonExistentArticleByID(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()

	req, err := http.NewRequest("GET", "/articles/100", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReq)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":404,"message":"Not Found","data":null}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

	tearDown()
}

func TestGetAllArticles(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReq)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"Success","data":[{"id":1,"title":"Post 1","content":"Some interesting content goes here.","author":"Farid"},{"id":2,"title":"Post 2","content":"Again some interesting content goes here.","author":"Farid"},{"id":3,"title":"Post 3","content":"Yet more interesting content goes here.","author":"Farid"}]}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

	tearDown()
}

func TestCreateArticle(t *testing.T) {
	setup()

	article.ConnectToDB(dataSourceTest)
	defer article.CloseDB()

	var jsonStr = []byte(`{"title":"New Article","content":"This is a new test article.","author":"Farid"}`)

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReq)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"status":201,"message":"Success","data":{"id":4}}`

	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

	tearDown()
}
