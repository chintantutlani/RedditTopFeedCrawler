package route

import (
	"net/http"

	"app/RedditTopFeedCrawler/handler"
	"app/RedditTopFeedCrawler/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(handler *handler.RouteHandlers) http.Handler {

	r := mux.NewRouter()
	m := &middleware.Middleware{}

	r.Use(m.Authentication)

	r.HandleFunc("/crawl", handler.PostsCrawlerHandler).Methods("POST")
	r.HandleFunc("/search", handler.SearchProfileHandler).Methods("GET")
	r.HandleFunc("/insights", handler.InsightsHandler).Methods("GET")

	routerWithAuth := m.Authentication(r)

	return routerWithAuth
}
