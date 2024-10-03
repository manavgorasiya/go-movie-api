package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Movie represents a movie with ID, Title, Director, Year, and Genre
type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
}

// In-memory movie store
var movies []Movie
var nextID int = 1

// CreateMovie handles the creation of a new movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	movie.ID = nextID
	nextID++
	movies = append(movies, movie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// GetMovies returns a list of all movies
func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// GetMovieByID retrieves a specific movie by its ID
func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	for _, movie := range movies {
		if movie.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// UpdateMovie updates an existing movie's details
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	for index, movie := range movies {
		if movie.ID == id {
			var updatedMovie Movie
			err := json.NewDecoder(r.Body).Decode(&updatedMovie)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			movies[index].Title = updatedMovie.Title
			movies[index].Director = updatedMovie.Director
			movies[index].Year = updatedMovie.Year
			movies[index].Genre = updatedMovie.Genre

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movies[index])
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// DeleteMovie deletes a movie by its ID
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	for index, movie := range movies {
		if movie.ID == id {
			// Remove movie from the slice
			movies = append(movies[:index], movies[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/movies", CreateMovie).Methods("POST")
	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", GetMovieByID).Methods("GET")
	router.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")

	// Start the server
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)
}
