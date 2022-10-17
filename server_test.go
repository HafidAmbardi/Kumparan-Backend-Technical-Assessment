package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Tests GET method without any query parameters
func TestGETArticle(t *testing.T) {
	t.Run("GET request to /article", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/article", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		response := httptest.NewRecorder()

		articleAPI(response, request)

		statusCode := response.Code
		body := response.Body.String()

		if statusCode != http.StatusOK {
			t.Errorf("got status %d want %d", statusCode, http.StatusOK)
		}
		if !IsJSON(body) {
			t.Errorf("Response body is not JSON, got %s", body)
		}
	})
}

// Tests POST method
func TestPostArticle(t *testing.T) {
	t.Run("POST request to /article", func(t *testing.T) {
		postArticle := article{"", "", "Dog", "Dog, the Wiki", "The dog is a domesticated descendant of the wolf...", time.Time{}}
		JSON, err := json.Marshal(postArticle)
		if err != nil {
			fmt.Println(err)
			return
		}

		request, err := http.NewRequest(http.MethodPost, "/article", bytes.NewReader(JSON))
		if err != nil {
			fmt.Println(err)
			return
		}

		response := httptest.NewRecorder()

		articleAPI(response, request)

		statusCode := response.Code
		body := response.Body.String()

		if statusCode != http.StatusCreated {
			t.Errorf("got status %d want %d", statusCode, http.StatusOK)
		}
		if !IsJSON(body) {
			t.Errorf("Response body is not JSON, got %s", body)
		}

		var responseAriticle article
		err = json.NewDecoder(strings.NewReader(body)).Decode(&responseAriticle)
		if err != nil {
			fmt.Println(err)
			return
		}

		if postArticle.AuthorName != responseAriticle.AuthorName || postArticle.Title != responseAriticle.Title || postArticle.Body != responseAriticle.Body {
			t.Errorf("Response body is not same as posted article, got %v want %v", responseAriticle, postArticle)
		}
		if responseAriticle.ID == "" {
			t.Errorf("ID is not generated")
		}
		if responseAriticle.AuthorID == "" {
			t.Errorf("AuthorID is not generated")
		}
		if responseAriticle.CreationTimeStamp.IsZero() {
			t.Errorf("CreatedAt is not generated")
		}
	})
}

// Tests GET method with "author" query parameter
func TestGETArticleAuthor(t *testing.T) {
	t.Run("GET request to /article?author=Dog", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/article?author=Dog", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		response := httptest.NewRecorder()

		articleAPI(response, request)

		statusCode := response.Code
		body := response.Body.String()

		var article []article
		json.Unmarshal([]byte(body), &article)

		if statusCode != http.StatusOK {
			t.Errorf("got status %d want %d", statusCode, http.StatusOK)
		}
		if !IsJSON(body) {
			t.Errorf("Response body is not JSON, got %s", body)
		}

		for _, element := range article {
			if element.AuthorName != "Dog" {
				t.Errorf("got status %s want Dog", element.AuthorName)
			}
		}
	})
}

// Tests GET method with "query" query parameter
func TestGETArticleQuery(t *testing.T) {
	t.Run("GET request to /article?query=wolf", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/article?query=wolf", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		response := httptest.NewRecorder()

		articleAPI(response, request)

		statusCode := response.Code
		body := response.Body.String()

		var nice []article
		json.Unmarshal([]byte(body), &nice)

		if statusCode != http.StatusOK {
			t.Errorf("got status %d want %d", statusCode, http.StatusOK)
		}
		if !IsJSON(body) {
			t.Errorf("Response body is not JSON, got %s", body)
		}

		for _, element := range nice {
			if !((strings.Contains(element.Title, "wolf")) || (strings.Contains(element.Body, "wolf"))) {
				t.Errorf("Title and Body does not contain wolf")
			}
		}
	})
}

// Checks if the format is JSON
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
