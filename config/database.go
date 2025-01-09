package config

import (
	"log"

	"github.com/gocql/gocql"
)

var SESSION *gocql.Session

func ConnectDB() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "bd2_proyecto1"
	cluster.Consistency = gocql.Quorum

	var err error
	if SESSION, err = cluster.CreateSession(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	InitDB()
}

func InitDB() {
	commands := []string{
		`CREATE TABLE IF NOT EXISTS users (
            user_id uuid PRIMARY KEY,
            name text,
            email text,
            password text
        );`,
		`CREATE TABLE IF NOT EXISTS movies (
            movie_id uuid PRIMARY KEY,
            poster_link text,
            series_title text,
            released_year text,
            certificate text,
            runtime text,
            genre text,
            IMDB_rating text,
            overview text,
            meta_score text,
            director text,
            star1 text,
            star2 text,
            star3 text,
            star4 text,
            no_of_votes text,
            gross text
        );`,
		`CREATE TABLE IF NOT EXISTS movies_watched_by_users (
            user_id uuid,
            movie_id uuid,
            watched_at timestamp,
            PRIMARY KEY ((user_id), movie_id)
        );`,
		`CREATE TABLE IF NOT EXISTS users_by_movie_watched (
            movie_id uuid,
            user_id uuid,
            watched_at timestamp,
            PRIMARY KEY ((movie_id), user_id)
        );`,
		`CREATE TABLE IF NOT EXISTS recommendations_by_user (
            user_id uuid,
            movie_id uuid,
            score float,
            PRIMARY KEY ((user_id), score, movie_id)
        ) WITH CLUSTERING ORDER BY (score DESC);`,
	}

	for _, command := range commands {
		if err := SESSION.Query(command).Exec(); err != nil {
			log.Fatalf("Error al ejecutar el comando: %v", err)
		}
	}
}
