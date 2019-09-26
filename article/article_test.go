package article

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

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

func TestNewArticle(t *testing.T) {
	setup()

	ConnectToDB(dataSourceTest)
	defer CloseDB()

	artcl, err := New("New Article", "Sample content for testing purpose.", "Developer")
	if err != nil {
		t.Fatal(err)
	}

	if *(artcl.ID) == int64(0) {
		err := errors.New("could not create new article")
		t.Fatal(err)
	}

	tearDown()
}

func TestFindByID(t *testing.T) {
	setup()

	ConnectToDB(dataSourceTest)
	defer CloseDB()

	artcl1, err := New("New Article to be found by id", "Sample content for testing purpose.", "Developer")
	if err != nil {
		t.Fatal(err)
	}

	artcl2, err := FindByID(*artcl1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if *artcl1.ID != *artcl2.ID {
		t.Fatal(err)
	}

	tearDown()
}

func TestFindNonExitentID(t *testing.T) {
	setup()

	ConnectToDB(dataSourceTest)
	defer CloseDB()

	_, err := FindByID(4)
	if err == nil {
		e := errors.New("error in finding non-existent article")
		t.Fatal(e)
	}

	tearDown()
}

func TestGetAll(t *testing.T) {
	setup()

	ConnectToDB(dataSourceTest)
	defer CloseDB()

	allArticles, err := GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(allArticles) != 3 {
		t.Fatal("Error in getting all articles.")
	}

	tearDown()
}
