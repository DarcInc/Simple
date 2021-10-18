package main

import (
	simple "Simple/model"
	"context"
	"github.com/darcinc/Simple/data"
	"github.com/darcinc/Simple/reflex"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
)

var imageRepo simple.ImageRepository

type SearchRequest struct {
}

type SearchResponse struct {
}

type ImageSearcher interface {
	Search(ctx context.Context, sr SearchRequest) (SearchResponse, error)
}

func NewImageSearcher(repository simple.ImageRepository) ImageSearcher {
	return nil
}

var imageSearcher ImageSearcher
var searchRepsonseTemplate template.Template

func ExtractSearchParamters(r *http.Request) SearchRequest {
	return SearchRequest{}
}

func findByLocationHandler(w http.ResponseWriter, r *http.Request) {
	request := ExtractSearchParamters(r)

	result, err := imageSearcher.Search(context.Background(), request)
	if err != nil {

	}

	searchRepsonseTemplate.ExecuteTemplate(w, "searchResponse", result)
}

func main() {
	DBURI := os.Getenv("DB_URI")

	poolConfig, err := pgxpool.ParseConfig(DBURI)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		os.Exit(1)


	}

	reflex := reflex.GlobalReflex()
	reflex.Register("caller", func(dm reflex.Reflex) (interface{}, bool) {
		return nil, true
	})


	caller := data.NewDBCaller(pool)              // 1 - application wide (stateless)
	server := data.NewMetadataServer(caller)      // 1 - application wide (stateless)
	imageRepo = simple.NewImageRepository(server) // 1 - application wide (stateless)
	imageSearcher = NewImageSearcher(imageRepo)
}
