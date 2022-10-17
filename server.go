package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// A struct to store the data of the articles
type article struct {
	ID                string    `json:"id"`
	AuthorID          string    `json:"author_id"`
	AuthorName        string    `json:"author_name"`
	Title             string    `json:"title"`
	Body              string    `json:"body"`
	CreationTimeStamp time.Time `json:"created_at"`
}

var DB, DBErr = sql.Open("mysql", "root:mewcat123@tcp(127.0.0.1:3306)/ricemilkDB2?parseTime=true")

// When a user makes a request to /article, this function will be called.
func articleAPI(w http.ResponseWriter, r *http.Request) {
	// Switch statement to check the RESTful method of the request.
	switch r.Method {

	// If the method is POST, we will add a new article to the DB.
	case "POST":
		// Deserializes the data from request's body to fit our "article" struct.
		var article article
		err := json.NewDecoder(r.Body).Decode(&article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Calls the postArticle function to store the article into the DB. The postArticle function returns an "article" struct with complete information (e.g., article ID, author ID, and creation time).
		article = postArticle(article)

		// Returns the article information back to the user as a response.
		w.Header().Set("content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(article)
		return

	case "GET":
		// Gets the optional query parameters from the request.
		queryQueryParam := r.URL.Query().Get("query")
		authorQueryParam := r.URL.Query().Get("author")

		// Gets the articles from the DB.
		result := getArticles(queryQueryParam, authorQueryParam)

		// Returns the articles from the DB to the user.
		w.Header().Set("content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return

	default:
		http.Error(w, "", http.StatusMethodNotAllowed) // Responds with 405 error if a different method is used.
	}
}

func postArticle(article article) article {
	// Generates a new article article ID and sets the creation time.
	var articleID string
	err := DB.QueryRow("SELECT id FROM articles ORDER by cast(id as unsigned) DESC LIMIT 1").Scan(&articleID)
	temp, err := strconv.Atoi(articleID)
	article.ID = strconv.Itoa(temp + 1)
	article.CreationTimeStamp = time.Now()

	// Searches for author ID using the author name
	var authorID string
	err = DB.QueryRow("SELECT id FROM authors WHERE name='" + article.AuthorName + "'").Scan(&authorID)

	// If the author name already has an assigned ID in the DB, use it.
	if authorID != "" {
		article.AuthorID = authorID
	} else { // If the author name doesn't have an assigned ID, generate a new one.
		err = DB.QueryRow("SELECT id FROM authors ORDER by cast(id as unsigned) DESC LIMIT 1").Scan(&authorID)
		temp, _ := strconv.Atoi(authorID)
		article.AuthorID = strconv.Itoa(temp + 1)

		// Stores the generated author ID in the "author" DB table.
		_, err := DB.Query("INSERT INTO authors VALUES('" + article.AuthorID + "','" + article.AuthorName + "')")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Stores the article in the "articles" DB table.
	stmtChild, err := DB.Prepare("INSERT INTO articles VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmtChild.Exec(article.ID, article.AuthorID, article.Title, article.Body, article.CreationTimeStamp)
	if err != nil {
		log.Fatal(err)
	}
	return article
}

func getArticles(query string, author string) []article {
	// Creates a string variable to store a query. The query will join "articles" table and "authors" table to retrieve articles using an author's name.
	sql := "SELECT articles.id, articles.author_id, articles.title, articles.body, articles.created_at, authors.name FROM articles INNER JOIN authors ON articles.author_id = authors.id"
	test := 0

	// If the user is using the "author" query parameter, find articles written by the author from the request.
	if author != "" {
		sql += (" WHERE authors.name='" + author + "'")
		test = 1
	}

	// If the user is using the "query" query parameter, find articles that contain the keyword from the request.
	if query != "" {
		if test != 1 {
			sql += " WHERE"
		} else {
			sql += " and"
		}
		sql += " (articles.title LIKE '%" + query + "%' or articles.body LIKE '%" + query + "%')"
	}

	sql += " ORDER by created_at DESC;"

	rows, err := DB.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	// Stores the result of the Query in an array of struct "article". Returns the array, therefore the information can be sent to the user.
	articles := []article{}
	for rows.Next() {
		var article article
		err := rows.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Body, &article.CreationTimeStamp, &article.AuthorName)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}

	return articles
}

func main() {
	if DBErr != nil {
		log.Fatal(DBErr)
	}
	defer DB.Close()

	http.HandleFunc("/article", articleAPI) // When a user makes a request to /article, call the articleAPI function.

	err := http.ListenAndServe("0.0.0.0:80", nil)
	log.Fatal(err)
}
