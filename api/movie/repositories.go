package movie

import (
	"errors"
	"example/api/schema"
	"example/config"

	"github.com/gocql/gocql"
)

func Paginate(query *gocql.Query, page int, pageSize int) ([]map[string]interface{}, error) {
	var nextPageToken []byte
	query = query.PageSize(pageSize)

	for i := 0; i < page-1; i++ {
		query = query.PageState(nextPageToken)
		if err := query.Exec(); err != nil {
			return nil, err
		}

		iter := query.Iter()
		for iter.MapScan(make(map[string]interface{})) {
		}
		nextPageToken = iter.PageState()

		iter.Close()
	}

	query = query.PageState(nextPageToken)
	iter := query.Iter()
	var results []map[string]interface{}
	var m = make(map[string]interface{})

	for {
		if !iter.MapScan(m) {
			break
		}
		results = append(results, m)
		m = make(map[string]interface{})
	}

	iter.Close()
	return results, nil
}

func ParseMovieModel(result map[string]interface{}) schema.MovieModel {
	return schema.MovieModel{
		MovieID:     result["movie_id"].(gocql.UUID),
		Poster:      result["poster_link"].(string),
		Title:       result["series_title"].(string),
		Released:    result["released_year"].(string),
		Certificate: result["certificate"].(string),
		Runtime:     result["runtime"].(string),
		Genre:       result["genre"].(string),
		Rating:      result["imdb_rating"].(string),
		Overview:    result["overview"].(string),
		Meta:        result["meta_score"].(string),
		Director:    result["director"].(string),
		Star1:       result["star1"].(string),
		Star2:       result["star2"].(string),
		Star3:       result["star3"].(string),
		Star4:       result["star4"].(string),
		Votes:       result["no_of_votes"].(string),
		Gross:       result["gross"].(string),
	}
}

func GetMoviesByFilterRepository(title string, genre string, page int, pageSize int) ([]schema.MovieModel, error) {
	session := config.SESSION
	var query *gocql.Query

	if title != "" && genre != "" {
		query = session.Query("SELECT * FROM movie WHERE series_title LIKE ? AND genre LIKE ? ALLOW FILTERING", "%"+title+"%", "%"+genre+"%")
	} else if title != "" {
		query = session.Query("SELECT * FROM movie WHERE series_title LIKE ?", "%"+title+"%")
	} else {
		query = session.Query("SELECT * FROM movie WHERE genre LIKE ?", "%"+genre+"%")
	}

	results, err := Paginate(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	var movies []schema.MovieModel
	for _, result := range results {
		movies = append(movies, ParseMovieModel(result))
	}

	return movies, nil
}

func GetMoviesRepository(page int, pageSize int) ([]schema.MovieModel, error) {
	session := config.SESSION

	query := session.Query("SELECT * FROM movie")
	results, err := Paginate(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	var movies []schema.MovieModel
	for _, result := range results {
		movies = append(movies, ParseMovieModel(result))
	}

	return movies, nil
}

func GetMovieByIdRepository(id string) (schema.MovieModel, error) {
	session := config.SESSION

	query := session.Query("SELECT * FROM movie WHERE movie_id = ?", id)
	model := make(map[string]interface{})
	if err := query.MapScan(model); err != nil {
		if err == gocql.ErrNotFound {
			return schema.MovieModel{}, errors.New("movie not found")
		}
		return schema.MovieModel{}, err
	}

	return ParseMovieModel(model), nil
}
