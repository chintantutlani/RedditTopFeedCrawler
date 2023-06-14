package server

import (
	"app/RedditTopFeedCrawler/dao"
	"app/RedditTopFeedCrawler/handler"
	"app/RedditTopFeedCrawler/service"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDAO() (*dao.DAO, error) {
	dbURL := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, err
	}

	return &dao.DAO{
		Client: client,
	}, nil
}

func InitService() (*service.Service, error) {
	// var service *service.Service

	dao, err := NewDAO()

	if err != nil {
		return nil, err
	}

	service := &service.Service{
		Dao: dao,
	}

	return service, nil

}

func InitRouteHandler() (*handler.RouteHandlers, error) {
	routeHandler := &handler.RouteHandlers{}

	service, err := InitService()
	if err != nil {
		return nil, err
	}

	routeHandler.Service = service

	return routeHandler, nil
}
