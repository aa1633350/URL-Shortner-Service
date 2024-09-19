package main

import (
	"UrlShortner/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var ctx = context.Background()

const baseUrl = "http://localhost:8000/"

// Request body structure for shortening URLs

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

func main() {
	r := chi.NewRouter()

	// Add some basic middleware using Chi
	r.Use(middleware.Logger)                    // Logs every request
	r.Use(middleware.Recoverer)                 // Recovers from panics
	r.Use(middleware.Timeout(60 * time.Second)) // Sets a request timeout

	// API Endpoints
	r.Post("/shorten", shortenURL)
	r.Get("/{shortenURL}", resolveURL)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

// handle post request to shorten URL
func shortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.LongURL == "" {
		http.Error(w, "Invalid input received", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode()

	client := database.CreateClient(0)
	defer client.Close()
	err = client.Set(ctx, shortCode, req.LongURL, 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	shortnedURL := baseUrl + shortCode
	response := map[string]string{"short_url": shortnedURL}
	w.Header().Set("Content-Type", "application/json")
	// converts map or struct to json, for example if below is the response then
	// response := map[string]string{
	//	"short_url": "http://localhost:8080/abc123",
	//}
	// then json.NewEncoder(w).Encode(response) will send the following JSON response to client
	//{
	//	"short_url": "http://localhost:8080/abc123"
	//}
	json.NewEncoder(w).Encode(response)

}

// resolveURL: Handles GET request to resolve the short URL to the original URL
func resolveURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")

	// retrieve long URL from redis
	client := database.CreateClient(0)
	defer client.Close()
	longURL, err := client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// redirect to original long URL
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)

}

// helper function to generate random short URL
func generateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	shortCode := make([]rune, 6)
	for i := range shortCode {
		shortCode[i] = chars[rand.Intn(len(chars))]
	}
	return string(shortCode)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
