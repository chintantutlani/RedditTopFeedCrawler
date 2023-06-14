package service

import (
	"app/RedditTopFeedCrawler/dao"
)

type Service struct {
	Dao *dao.DAO
}

func (cd *Service) CrawlData(dbName string, collectionName string) error {

	err := cd.Dao.CrawlPosts(dbName, collectionName)

	return err

}
