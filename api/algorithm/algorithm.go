package algorithm

import (
	"example/config"
	"fmt"
	"log"
	"math"
	"sort"
)

type PreResult struct {
	movieWatched    string
	movieNotWatched string
	score           float64
}

type Result struct {
	movieNotWatched string
	score           float64
}

func GenerateRecommendation(user_id string) ([]Result, error) {
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
		return nil, err
	}

	if len(watchedMoviesId) == 0 {
		return nil, nil
	}

	// Fetch movies watched by anyone
	iter = session.Query("SELECT DISTINCT movie_id FROM user_by_movie_watched").Iter()
	var movieId string
	for iter.Scan(&movieId) {
		allMoviesId = append(allMoviesId, movieId)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	// Fetch all users that has watched a movie
	iter = session.Query("SELECT DISTINCT user_id FROM movie_watched_by_user").Iter()
	var userID string
	for iter.Scan(&userID) {
		allUsersId = append(allUsersId, userID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	fmt.Println("allUsersId: ", len(allUsersId), "allMoviesId: ", len(allMoviesId), "watchedMoviesId: ", len(watchedMoviesId))

	result, err := Algorithm(user_id, watchedMoviesId, allUsersId, allMoviesId)
	if err != nil {
		return nil, err
	}
	return result, nil

}

// Using recommendation algorithm Hybrid (HeatS and ProbS) to create recommendations
// based on the movies watched by the user and the users who watched them
func Algorithm(user_id string, moviesWatched []string, usersId []string, moviesId []string) ([]Result, error) {
	lambda := 0.5
	var results []PreResult
	log.Println("STARTED")

	for _, movieWatched := range moviesWatched { // For each movie watched
		degreeMovieWatched, err := calculateDegreeOfMovie(movieWatched) // A
		if err != nil {
			return nil, err
		}

		for i, movieId := range moviesId { // For each movie in the database
			if checkIfMovieIsWatched(moviesWatched, movieId) { // if the movie is watched by the user
				continue
			}

			log.Println("i: ", i)
			log.Println("movieWatched: ", movieWatched, "movieId: ", movieId)

			degreeMovieNotWatched, err := calculateDegreeOfMovie(movieId) // B
			if err != nil {
				return nil, err
			}
			if degreeMovieNotWatched == 0 { // if the movie is not watched by anyone
				log.Println("movie not watched is not watched by anyone")
				continue
			}

			leftPart := leftPart(degreeMovieNotWatched, degreeMovieWatched, lambda)
			log.Println("leftPart: ", leftPart)
			rightPart, err := rightPart(usersId, movieWatched, movieId)
			log.Println("rightPart: ", rightPart)
			if err != nil {
				return nil, err
			}
			score := leftPart * rightPart
			log.Println("score: ", score)

			results = append(results, PreResult{movieWatched, movieId, score})
		}
	}

	return GenerateFinalScore(results), nil
}

func leftPart(degreeMovieNotWatched int, degreeMovieWatched int, lambda float64) float64 {
	return 1 / (math.Pow(float64(degreeMovieWatched), (1-lambda)) * math.Pow(float64(degreeMovieNotWatched), (lambda)))
}

func rightPart(usersId []string, movieWatchedId string, movieNotWatched string) (float64, error) {
	var result float64
	usersMovieCount, err := getUsersMovieCount(usersId)
	log.Println("USERS MOVIE COUNT DONE")
	if err != nil {
		return 0, err
	}

	for _, userId := range usersId {
		// check the number of movies which the user has interacted
		kj := usersMovieCount[userId]
		// avoid division by zero
		if kj == 0 {
			continue
		}

		a1 := 0
		a2 := 0

		if watched, err := checkIfUserWatchedTheMovie(userId, movieWatchedId); err != nil {
			if watched {
				a1 = 1
			}
		}
		if watched, err := checkIfUserWatchedTheMovie(userId, movieNotWatched); err != nil {
			if watched {
				a2 = 1
			}
		}

		result += (float64(a1) * float64(a2)) / float64(kj)
	}

	return result, nil
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

func GenerateFinalScore(results []PreResult) []Result {

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

	var finalResultsSlice []Result
	for movie, score := range finalResults {
		finalResultsSlice = append(finalResultsSlice, Result{movie, score})
	}

	sort.Slice(finalResultsSlice, func(i, j int) bool {
		return finalResultsSlice[i].score > finalResultsSlice[j].score
	})

	return finalResultsSlice
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
