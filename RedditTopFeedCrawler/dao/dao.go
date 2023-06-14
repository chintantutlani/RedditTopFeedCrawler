package dao

import (
	models "app/RedditTopFeedCrawler/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type DAO struct {
	Client *mongo.Client
}

func (da *DAO) CrawlPosts(dbName, collectionName string) error {
	db := da.Client.Database(dbName)
	collection := db.Collection(collectionName)

	url := "https://www.reddit.com/top/.json"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return errors.New("Invalid response format: 'data' field not found")
	}

	children, ok := data["children"].([]interface{})
	if !ok {
		return errors.New("Invalid response format: 'children' field not found")
	}

	if len(children) == 0 {
		return errors.New("No posts found in the response.")
	}

	var documents []interface{}
	for _, child := range children {
		doc, ok := child.(map[string]interface{})["data"].(map[string]interface{})
		if ok {
			documents = append(documents, doc)
		}
	}

	if len(documents) == 0 {
		return errors.New("No valid posts found in the response.")
	}

	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	for _, child := range children {
		postData, ok := child.(map[string]interface{})["data"].(map[string]interface{})
		if ok {
			postTitle, _ := postData["title"].(string)
			_, err := collection.InsertOne(context.Background(), postData)
			if err != nil {
				log.Println("Failed to insert post:", postTitle)
			} else {
				log.Println("Inserted post:", postTitle)
			}
		}
	}

	return nil
}

func (dao *DAO) SearchProfile(w http.ResponseWriter, r *http.Request) {
	subreddits, ok := r.URL.Query()["subreddit"]
	if !ok || len(subreddits) == 0 {
		http.Error(w, "Missing required query parameter 'subreddit'", http.StatusBadRequest)
		return
	}

	collection := dao.Client.Database("reddit").Collection("nn")
	posts := make([]models.Post, 0)

	for _, subreddit := range subreddits {
		filter := bson.M{"subreddit": subreddit}
		cur, err := collection.Find(context.Background(), filter)
		if err != nil {
			log.Fatal("Failed to retrieve posts:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			var post models.Post
			err := cur.Decode(&post)
			if err != nil {
				log.Fatal("Failed to decode post:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError)
		return
	}
}

func (dao *DAO) GetInsights(dbName, collectionName string) (*models.InsightsResponse, error) {
	db := dao.Client.Database(dbName)
	collection := db.Collection(collectionName)

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$subreddit",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"count": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var totalPosts int
	var insights [4]int

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		count, ok := result["count"].(int32)
		if !ok {
			return nil, err
		}

		totalPosts += int(count)
		if count < 100 {
			insights[0]++
		} else if count <= 10000 {
			insights[1]++
		} else if count <= 100000 {
			insights[2]++
		} else {
			insights[3]++
		}
	}

	pipeline = []bson.M{
		{
			"$sort": bson.M{
				"ups": -1,
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err = collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var topPosts []models.Post

	for cursor.Next(context.Background()) {
		var post models.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		topPosts = append(topPosts, post)
	}

	response := &models.InsightsResponse{
		TotalPosts:           totalPosts,
		PostsLessThan100:     insights[0],
		Posts101To10000:      insights[1],
		Posts10001To100000:   insights[2],
		PostsGreaterThan100K: insights[3],
		Timestamp:            time.Now().Unix(),
		TopPosts:             topPosts,
	}

	return response, nil
}

func (da *DAO) DaoCrawler() {
	dbName := "reddit"
	collectionName := "mycollection"

	err := da.CrawlPosts(dbName, collectionName)
	if err != nil {
		log.Fatal("Failed to crawl posts:", err)
	}
}
