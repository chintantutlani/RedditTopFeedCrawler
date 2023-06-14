package handler

import (
	// "app/RedditTopFeedCrawler/dao"

	"app/RedditTopFeedCrawler/service"
	"encoding/json"

	"log"
	"net/http"
)

type RouteHandlers struct {
	Service *service.Service
}

// type PostCrawler struct {
// 	dao *dao.DAO
// }

func (rh *RouteHandlers) PostsCrawlerHandler(w http.ResponseWriter, r *http.Request) {

	dbName := "reddit"
	collectionName := "profile"
	err := rh.Service.CrawlData(dbName, collectionName)
	if err != nil {
		log.Fatal("Failed to crawl posts:", err)
		http.Error(w, "Internal Servker Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Posts crawled and stored in MongoDB",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Failed to encode response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseBytes)
}

func (dao *RouteHandlers) SearchProfileHandler(w http.ResponseWriter, r *http.Request) {
	dao.Service.Dao.SearchProfile(w, r)
}

func (s *RouteHandlers) InsightsHandler(w http.ResponseWriter, r *http.Request) {
	dbName := "reddit"
	collectionName := "profile"
	response, err := s.Service.Dao.GetInsights(dbName, collectionName)
	if err != nil {
		log.Fatal("Failed to retrieve insights:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Failed to encode response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseBytes)
}
