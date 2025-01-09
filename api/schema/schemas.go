package schema

import (
	"github.com/gocql/gocql"
)

type Movie struct {
	MovieID     gocql.UUID `json:"movie_id"`
	Poster      string     `json:"poster"`
	Title       string     `json:"title"`
	Released    string     `json:"released"`
	Certificate string     `json:"certificate"`
	Runtime     string     `json:"runtime"`
	Genre       string     `json:"genre"`
	Rating      string     `json:"rating"`
	Overview    string     `json:"overview"`
	Meta        string     `json:"meta"`
	Director    string     `json:"director"`
	Star1       string     `json:"star1"`
	Star2       string     `json:"star2"`
	Star3       string     `json:"star3"`
	Star4       string     `json:"star4"`
	Votes       string     `json:"votes"`
	Gross       string     `json:"gross"`
}
