package schema

import (
	"github.com/gocql/gocql"
)

type MovieModel struct {
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

type UserRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    gocql.UUID `json:"user_id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
}

type UserUpdate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RecommendationCreate struct {
	UserID  string  `json:"user_id"`
	MovieID string  `json:"movie_id"`
	Score   float64 `json:"score"`
}

type RecommendationResponse struct {
	MovieID string  `json:"movie_id"`
	Score   float32 `json:"score"`
}
