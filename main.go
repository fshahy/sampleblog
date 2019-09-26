package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fshahy/sampleblog/article"
	"github.com/fshahy/sampleblog/helper"
	_ "github.com/lib/pq"
)

var dataSource = fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=disable",
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PASSWORD"),
	os.Getenv("POSTGRES_DB"),
)

func main() {

	article.ConnectToDB(dataSource)
	defer article.CloseDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleReq)

	log.Fatal(http.ListenAndServe(":8080", mux))

}

// handle all routes
func handleReq(w http.ResponseWriter, r *http.Request) {

	// drop extra slashes
	tmpPath := strings.TrimPrefix(r.URL.String(), "/")
	reqPath := strings.TrimSuffix(tmpPath, "/")

	// extract url path and parameter (if any)
	pathTokens := strings.Split(reqPath, "/")

	if len(pathTokens) == 1 {
		if pathTokens[0] == "articles" {
			if r.Method == http.MethodGet {
				var data []article.Article
				allArticles, err := article.GetAll()
				if err != nil {
					log.Println(err)
				}

				if len(allArticles) > 0 {
					for _, artcl := range allArticles {
						data = append(data, artcl)
					}

					writeResponse(w, 200, "Success", data)
					return
				}
				writeResponse(w, 404, "Not Found", []interface{}{})
				return
			} else if r.Method == http.MethodPost {
				decoder := json.NewDecoder(r.Body)
				var a article.Article
				err := decoder.Decode(&a)
				if err != nil {
					log.Println(err)
				}

				// validate posted data
				isValid := validateArticleData(a)

				if isValid {
					artcl, err := article.New(*a.Title, *a.Content, *a.Author)
					if err != nil {
						log.Println(err)
					}

					writeResponse(w, 201, "Success", *artcl.ID)
					return
				}
				writeResponse(w, 406, "Not Acceptable", []interface{}{})
				return
			} else {
				writeResponse(w, 405, "Method Not Allowed", []interface{}{})
				return
			}
		} else {
			writeResponse(w, 404, "Not Found", []interface{}{})
			return
		}
	} else if len(pathTokens) == 2 {
		if pathTokens[0] == "articles" {
			if r.Method == http.MethodGet {
				var artcl article.Article

				resourceID, err := strconv.ParseInt(pathTokens[1], 10, 64)
				if err != nil {
					log.Println(err)
				}

				artcl, err = article.FindByID(resourceID)
				if err != nil {
					writeResponse(w, 404, "Not Found", []interface{}{})
					return
				}

				writeResponse(w, 200, "Success", []article.Article{artcl})
				return
			}
			writeResponse(w, 405, "Method Not Allowed", []interface{}{})
			return
		}
		writeResponse(w, 404, "Not Found", []interface{}{})
		return
	} else {
		writeResponse(w, 404, "Not Found", []interface{}{})
		return
	}
}

func validateArticleData(artcl article.Article) bool {
	if *artcl.Title == "" || *artcl.Content == "" || *artcl.Author == "" {
		return false
	}

	// TODO: add some business related validations

	return true
}

func writeResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	// fmt.Println(status)
	response, err := helper.NewResponse(status, message, data)
	if err != nil {
		log.Println(err)
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)
}
