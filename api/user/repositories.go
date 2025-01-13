package user

import (
	"example/api/schema"
	"example/config"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

func GetUserRecommendationsRepository(user_id string) ([]schema.RecommendationResponse, error) {
	session := config.SESSION
	var recommendations []schema.RecommendationResponse
	var rec schema.RecommendationResponse

	iter := session.Query(`SELECT movie_id, score FROM recommendation_by_user WHERE user_id = ?`, user_id).Iter()

	for iter.Scan(&rec.MovieID, &rec.Score) {
		recommendations = append(recommendations, rec)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return recommendations, nil
}

func GetUserHistoryRepository(user_id string) ([]string, error) {
	session := config.SESSION
	var movies_id []string
	var movie_id string

	if err := session.Query("SELECT user_id FROM user_by_id WHERE user_id = ?", user_id).Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	iter := session.Query(`SELECT movie_id FROM movie_watched_by_user WHERE user_id = ?`, user_id).Iter()
	for iter.Scan(&movie_id) {
		movies_id = append(movies_id, movie_id)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return movies_id, nil
}

func WatchMovieRepository(user_id string, movie_id string) error {
	session := config.SESSION

	if err := session.Query("SELECT user_id FROM user_by_id WHERE user_id = ?", user_id).Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	if err := session.Query("SELECT movie_id FROM movie WHERE movie_id = ?", movie_id).Exec(); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("movie not found")
		}
		return err
	}

	if err := session.Query(`INSERT INTO movie_watched_by_user (user_id, movie_id, watched_at) VALUES (?, ?, ?)`, user_id, movie_id, time.Now()).Exec(); err != nil {
		return err
	}

	if err := session.Query(`INSERT INTO user_by_movie_watched (user_id, movie_id, watched_at) VALUES (?, ?, ?)`, user_id, movie_id, time.Now()).Exec(); err != nil {
		return err
	}

	if err := session.Query(`DELETE FROM recommendation_by_user WHERE user_id = ? AND movie_id = ?`, user_id, movie_id).Exec(); err != nil {
		return err
	}

	return nil
}

func GetUsersRepository() ([]schema.UserResponse, error) {
	var users []schema.UserResponse
	session := config.SESSION
	iter := session.Query(`SELECT user_id, name, email FROM user_by_id`).Iter()
	var user schema.UserResponse
	for iter.Scan(&user.ID, &user.Name, &user.Email) {
		users = append(users, user)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByRepository(id string) (schema.UserResponse, error) {
	var user schema.UserResponse
	session := config.SESSION
	if err := session.Query(`SELECT user_id, name, email FROM user_by_id WHERE user_id = ?`, id).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return user, err
	}
	return user, nil
}

func UpdateUserRepository(id string, user schema.UserUpdate) error {
	session := config.SESSION
	var existingEmail string
	if err := session.Query(`SELECT email FROM user_by_id WHERE user_id = ? LIMIT 1`, id).Scan(&existingEmail); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	if err := session.Query(`UPDATE user_by_id SET name = ?, email = ?, password = ? WHERE user_id = ?`, user.Name, user.Email, user.Password, id).Exec(); err != nil {
		return err
	}

	if existingEmail != user.Email {
		if err := session.Query(`DELETE FROM user_by_email WHERE email = ?`, existingEmail).Exec(); err != nil {
			return err
		}
		if err := session.Query(`INSERT INTO user_by_email (email, user_id, name, password) VALUES (?, ?, ?, ?)`, user.Email, id, user.Name, user.Password).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func DeleteUserRepository(id string) error {
	session := config.SESSION
	var existingEmail string
	if err := session.Query(`SELECT email FROM user_by_id WHERE user_id = ? LIMIT 1`, id).Scan(&existingEmail); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	if err := session.Query(`DELETE FROM user_by_id WHERE user_id = ?`, id).Exec(); err != nil {
		return err
	}
	if err := session.Query(`DELETE FROM user_by_email WHERE email = ?`, existingEmail).Exec(); err != nil {
		return err
	}
	return nil
}
