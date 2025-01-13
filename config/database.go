package config

import (
	"log"
	"os"

	"github.com/gocql/gocql"
)

var SESSION *gocql.Session

func ConnectDB() {
	host := os.Getenv("CASSANDRA_HOST")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	var err error
	if SESSION, err = cluster.CreateSession(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	InitDB()
}

func InitDB() {
	commands := []string{
		`CREATE TABLE IF NOT EXISTS movie (
            movie_id uuid,
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
            gross text,
            PRIMARY KEY (movie_id)
        );`,

		`CREATE CUSTOM INDEX IF NOT EXISTS idx_title ON movie (series_title) 
            USING 'org.apache.cassandra.index.sasi.SASIIndex' 
            WITH OPTIONS = {
            'mode': 'CONTAINS', 
            'analyzer_class': 'org.apache.cassandra.index.sasi.analyzer.NonTokenizingAnalyzer', 
            'case_sensitive': 'false'};`,

		`CREATE CUSTOM INDEX IF NOT EXISTS idx_genre ON movie (genre) 
            USING 'org.apache.cassandra.index.sasi.SASIIndex' 
            WITH OPTIONS = {
            'mode': 'CONTAINS', 
            'analyzer_class': 'org.apache.cassandra.index.sasi.analyzer.NonTokenizingAnalyzer', 
            'case_sensitive': 'false'};`,

		`CREATE TABLE IF NOT EXISTS movie_watched_by_user (
            user_id uuid,
            movie_id uuid,
            watched_at timestamp,
            PRIMARY KEY ((user_id), movie_id)
        );`,

		`CREATE TABLE IF NOT EXISTS user_by_movie_watched (
            movie_id uuid,
            user_id uuid,
            watched_at timestamp,
            PRIMARY KEY ((movie_id), user_id)
        );`,

		`CREATE TABLE IF NOT EXISTS recommendation_by_user (
            user_id uuid,
            movie_id uuid,
            score float,
            PRIMARY KEY ((user_id), movie_id)
        )`,

		`CREATE TABLE IF NOT EXISTS user_by_id (
            user_id uuid,
            name text,
            email text,
            password text,
            PRIMARY KEY (user_id)
        );`,

		`CREATE TABLE IF NOT EXISTS user_by_email (
            email text,
            user_id uuid,
            name text,
            password text,
            PRIMARY KEY (email)
        );`,
	}

	for _, command := range commands {
		if err := SESSION.Query(command).Exec(); err != nil {
			log.Fatalf("Error al ejecutar el comando: %v", err)
		}
	}
}
