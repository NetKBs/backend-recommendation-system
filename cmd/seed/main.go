package main

import (
	"encoding/json"
	"example/api/algorithm"
	"example/api/schema"
	"example/api/user"
	"example/config"
	"io"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	defer config.SESSION.Close()
	LoadData()
}

func LoadData() {
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
	var moviesId []string

	for _, movie := range movies {
		movie.MovieID = gocql.UUID(uuid.New())
		moviesId = append(moviesId, movie.MovieID.String())

		if err := config.SESSION.Query(
			`INSERT INTO movie (movie_id, poster_link, series_title, released_year, certificate, runtime, genre, IMDB_rating, overview, meta_score, director, star1, star2, star3, star4, no_of_votes, gross) 
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
			log.Fatalf("Error al insertar la pelicula en movie_by_id: %v", err)
		}
	}

	log.Println("Movies inserted")

	// Users
	gofakeit.Seed(111)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error al generar la contraseña: %v", err)
	}

	var usersId []string

	for i := 0; i < 30; i++ {
		newUUID := gocql.UUID(uuid.New())
		user := schema.UserRegister{
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: string(hashedPassword),
		}

		if err := config.SESSION.Query(`INSERT INTO user_by_email (email, user_id, name, password) VALUES (?, ?, ?, ?)`,
			user.Email, newUUID, user.Name, user.Password).Exec(); err != nil {
			log.Fatalf("Error al insertar el usuario en user_by_email: %v", err)
		}
		if err := config.SESSION.Query(`INSERT INTO user_by_id (user_id, name, email, password) VALUES (?, ?, ?, ?)`,
			newUUID, user.Name, user.Email, user.Password).Exec(); err != nil {
			log.Fatalf("Error al insertar el usuario en user_by_id: %v", err)
		}

		usersId = append(usersId, newUUID.String())
		log.Println("User ", i, " inserted")
	}

	// Watched movies
	for _, user_id := range usersId {
		numMoviesToWatch := gofakeit.Number(1, 100)

		for i := 0; i < numMoviesToWatch; i++ {
			movie_id := moviesId[i]

			if err := user.WatchMovieRepository(user_id, movie_id); err != nil {
				log.Fatalf("Error al insertar la pelicula en movie_watched_by_user: %v", err)
			}

		}

	}
	log.Println("Movies watched inserted")

	// Generate new recommendations
	for i, user_id := range usersId {
		algorithm.GenerateRecommendation(user_id)
		log.Print("Recommendations generated for user ", i)
	}

}
