package main

import (
	"encoding/json"
	"example/api/schema"
	"example/config"
	"io"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func main() {
	config.ConnectDB()
	defer config.SESSION.Close()
	LoadMovies()
}

func LoadMovies() {
	f, err := os.Open("movie_data.json")
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
	}

	var movies []schema.MovieModel
	if err := json.Unmarshal(bytes, &movies); err != nil {
		log.Fatalf("Error al deserializar el archivo: %v", err)
	}

	for _, movie := range movies {
		movie.MovieID = gocql.UUID(uuid.New())
		if err := config.SESSION.Query(
			`INSERT INTO movies (movie_id, poster_link, series_title, released_year, certificate, runtime, genre, IMDB_rating, overview, meta_score, director, star1, star2, star3, star4, no_of_votes, gross) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			movie.MovieID,
			movie.Poster,
			movie.Title,
			movie.Released,
			movie.Certificate,
			movie.Runtime,
			movie.Genre,
			movie.Rating,
			movie.Overview,
			movie.Meta,
			movie.Director,
			movie.Star1,
			movie.Star2,
			movie.Star3,
			movie.Star4,
			movie.Votes,
			movie.Gross,
		).Exec(); err != nil {
			log.Fatalf("Error al insertar la pelicula: %v", err)
		}
	}

}
