package model

type Post struct {
	ID        string `bson:"id"`
	Title     string `bson:"title"`
	Name      string `bson:"name"`
	Ups       int    `bson:"ups"`
	Subreddit string `bson:"subreddit"`
}

// type Response struct {
// 	Data struct {
// 		Children []struct {
// 			Data Post `json:"data"`
// 		} `json:"children"`
// 	} `json:"data"`
// }

type InsightsResponse struct {
	TotalPosts           int    `json:"totalPosts"`
	PostsLessThan100     int    `json:"postsLessThan100"`
	Posts101To10000      int    `json:"posts101To10000"`
	Posts10001To100000   int    `json:"posts10001To100000"`
	PostsGreaterThan100K int    `json:"postsGreaterThan100K"`
	Timestamp            int64  `json:"timestamp"`
	TopPosts             []Post `json:"topPosts"`
}
