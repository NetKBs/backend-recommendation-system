package algorithm

import (
	"example/api/schema"
	"example/config"
	"fmt"
	"log"
	"math"
	"sync"
)

type PreResult struct {
	movieWatched    string
	movieNotWatched string
	score           float64
}

func GenerateRecommendation(user_id string) {
	session := config.SESSION
	var allMoviesId []string
	var allUsersId []string
	var watchedMoviesId []string

	// Fetch movies watched by the user
	iter := session.Query(`SELECT movie_id FROM movie_watched_by_user WHERE user_id = ?`, user_id).Iter()
	var movieID string
	for iter.Scan(&movieID) {
		watchedMoviesId = append(watchedMoviesId, movieID)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	if len(watchedMoviesId) == 0 {
		log.Println("No movies watched by the user")
	}

	// Fetch movies watched by anyone
	iter = session.Query("SELECT DISTINCT movie_id FROM user_by_movie_watched").Iter()
	var movieId string
	for iter.Scan(&movieId) {
		allMoviesId = append(allMoviesId, movieId)
	}
	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	// Fetch all users that has watched a movie
	iter = session.Query("SELECT DISTINCT user_id FROM movie_watched_by_user").Iter()
	var userID string
	for iter.Scan(&userID) {
		allUsersId = append(allUsersId, userID)
	}
	if err := iter.Close(); err != nil {
		log.Println(err)
	}

	fmt.Println("allUsersId: ", len(allUsersId), "allMoviesId: ", len(allMoviesId), "watchedMoviesId: ", len(watchedMoviesId))

	results, err := Algorithm(user_id, watchedMoviesId, allUsersId, allMoviesId)
	if err != nil {
		log.Println(err)
	}

	// save
	for _, result := range results {
		score := float32(result.Score)
		if err := config.SESSION.Query(`INSERT INTO recommendation_by_user (user_id, movie_id, score) VALUES (?, ?, ?)`, result.UserID, result.MovieID, score).Exec(); err != nil {
			log.Printf("Error: %v\n", err)
		}
	}

	log.Println("Algorithm finished")
}

func Algorithm(user_id string, moviesWatched []string, usersId []string, moviesId []string) ([]schema.RecommendationCreate, error) {
	lambda := 0.5
	var results []PreResult
	var mu sync.Mutex     // se usa para proteger el acceso concurrente a variables compartidas
	var wg sync.WaitGroup // se para esperar que un conjunto de goroutines se complete

	usersMovieCount, err := getUsersMovieCount(usersId)
	if err != nil {
		return nil, err
	}

	for _, movieWatched := range moviesWatched { // For each movie watched
		degreeMovieWatched, err := calculateDegreeOfMovie(movieWatched) // A
		if err != nil {
			return nil, err
		}

		for _, movieId := range moviesId { // For each movie in the database
			if checkIfMovieIsWatched(moviesWatched, movieId) { // if the movie is watched by the user
				continue
			}

			wg.Add(1)

			go func(movieWatched, movieId string) {
				defer wg.Done()

				degreeMovieNotWatched, err := calculateDegreeOfMovie(movieId) // B
				if err != nil {
					return
				}

				leftPart := leftPart(degreeMovieNotWatched, degreeMovieWatched, lambda)
				rightPart, err := rightPart(usersId, movieWatched, movieId, usersMovieCount)
				if err != nil {
					return
				}
				score := leftPart * rightPart

				mu.Lock()
				results = append(results, PreResult{movieWatched, movieId, score})
				mu.Unlock()

			}(movieWatched, movieId)

		}
	}

	wg.Wait()
	finalResults, err := GenerateFinalScore(results, user_id)
	if err != nil {
		return nil, err
	}

	return finalResults, nil
}

func leftPart(degreeMovieNotWatched int, degreeMovieWatched int, lambda float64) float64 {
	return 1 / (math.Pow(float64(degreeMovieWatched), (1-lambda)) * math.Pow(float64(degreeMovieNotWatched), (lambda)))
}

func rightPart(usersId []string, movieWatchedId string, movieNotWatched string, usersMovieCount map[string]int) (float64, error) {
	var result float64

	for _, userId := range usersId {
		// check the number of movies which the user has interacted
		kj := usersMovieCount[userId]

		a1 := 0
		a2 := 0

		watched, err := checkIfUserWatchedTheMovie(userId, movieWatchedId)
		if err != nil {
			return 0, err
		}
		if watched {
			a1 = 1
		}

		watched, err = checkIfUserWatchedTheMovie(userId, movieNotWatched)
		if err != nil {
			return 0, err
		}
		if watched {
			a2 = 1
		}

		result += (float64(a1) * float64(a2)) / float64(kj)
	}

	return result, nil
}

func GenerateFinalScore(results []PreResult, userId string) ([]schema.RecommendationCreate, error) {

	// Normalize the scores
	var total float64
	for _, result := range results {
		total += result.score
	}
	for _, result := range results {
		result.score = result.score / total
	}

	// Aggregate the scores
	var finalResults = make(map[string]float64)
	for _, result := range results {
		finalResults[result.movieNotWatched] += result.score
	}

	var finalResultsSlice []schema.RecommendationCreate

	for movieId, score := range finalResults {
		result := schema.RecommendationCreate{UserID: userId, MovieID: movieId, Score: score}
		finalResultsSlice = append(finalResultsSlice, result)
	}

	return finalResultsSlice, nil
}

func getUsersMovieCount(usersId []string) (map[string]int, error) {
	session := config.SESSION
	userMovieCount := make(map[string]int)

	for _, userId := range usersId {
		var count int
		if err := session.Query(`SELECT COUNT(movie_id) FROM movie_watched_by_user WHERE user_id = ?`, userId).Scan(&count); err != nil {
			return nil, err
		}
		userMovieCount[userId] = count
	}

	return userMovieCount, nil
}

func calculateDegreeOfMovie(movieId string) (int, error) {
	var degree int
	if err := config.SESSION.Query(`SELECT COUNT(user_id) FROM user_by_movie_watched WHERE movie_id = ?`, movieId).Scan(&degree); err != nil {
		return 0, err
	}
	return degree, nil
}

func checkIfUserWatchedTheMovie(userId string, movieId string) (bool, error) {
	var count int
	if err := config.SESSION.Query(`SELECT COUNT(movie_id) FROM movie_watched_by_user WHERE user_id = ? AND movie_id = ?`, userId, movieId).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func checkIfMovieIsWatched(moviesWatched []string, movieId string) bool {
	for _, movieWatched := range moviesWatched {
		if movieWatched == movieId {
			return true
		}
	}
	return false
}
