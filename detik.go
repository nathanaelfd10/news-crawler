package main

import (
	"context"
	c "detik-scraper/config"
	db "detik-scraper/database"
	m "detik-scraper/models"
	u "detik-scraper/utils"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	fmt.Println(c.MongoDBHost)
	mongoDBURI := fmt.Sprintf("mongodb://%s:%s", c.MongoDBHost, c.MongoDBPort)
	fmt.Println(mongoDBURI)
	client, collection, err := db.ConnectToMongoDB(mongoDBURI, c.DatabaseName, c.CollectionNameDetik)
	if err != nil {
		log.Fatalf("Error establishing MongoDB connection: %v", err)
	}
	defer client.Disconnect(context.TODO())

	c := colly.NewCollector()

	maxPages := 1

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawling", r.URL.String())
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		article := extractDetikArticle(e)
		fmt.Println("Found: ", article.ArticlePreview.Title)

		content, err := GetContent(article.ArticlePreview.URL)
		if err != nil {
			fmt.Errorf("Can't fetch content, %v", err)
		}

		article.Content = content

		_, err = collection.InsertOne(context.TODO(), article)
		if err != nil {
			log.Printf("\tFailed to insert article: %s", err)
			return
		}
	})

	for page := 1; page <= maxPages; page++ {
		url := fmt.Sprintf("https://news.detik.com/indeks/%d", page)

		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting page %d: %s\n", page, err)
			break
		}
	}

	fmt.Println("Crawling completed.")
}

func extractDetikArticle(e *colly.HTMLElement) m.ArticleDetik {
	thumbnailURL := e.ChildAttr(".media__image img", "src")
	thumbnail, err := u.GetImage(thumbnailURL)
	if err != nil {
		fmt.Errorf("Failed to fetch image: %s", err)
	}

	article := m.ArticleDetik{
		ArticlePreview: m.ArticlePreview[m.DetikSpecificInfo]{
			Title:        e.ChildText(".media__title > a"),
			URL:          e.ChildAttr(".media__title > a", "href"),
			Thumbnail:    thumbnail,
			ThumbnailURL: thumbnailURL,
			WebSpecificInfo: m.DetikSpecificInfo{
				TimestampUTC: e.ChildAttr(".media__date > span", "d-time"),
				DateWIB:      e.ChildAttr(".media__date > span", "title"),
			},
			CrawledAt: time.Now(),
		},
	}

	// article := m.ArticleDetik{
	// 	Title:        e.ChildText(".media__title > a"),
	// 	URL:          e.ChildAttr(".media__title > a", "href"),
	// 	Thumbnail:    thumbnail,
	// 	ThumbnailURL: thumbnailURL,
	// }

	// article := m.ArticleDetik{
	// 	BaseInfo: m.BaseInfo{
	// 		Title:        e.ChildText(".media__title > a"),
	// 		URL:          e.ChildAttr(".media__title > a", "href"),
	// 		Thumbnail:    thumbnail,
	// 		ThumbnailURL: thumbnailURL,
	// 		CrawledAt: ,
	// 	},
	// 	TimestampUTC: e.ChildAttr(".media__date > span", "d-time"),
	// 	DateWIB:      e.ChildAttr(".media__date > span", "title"),
	// }

	return article
}

func GetContent(url string) (m.Content, error) {
	var content m.Content

	c := colly.NewCollector()

	c.OnHTML(".content__bg", func(e *colly.HTMLElement) {
		e.DOM.Find(".staticdetail_container, .lihatjg, .para_caption, .staticdetail_ads").Remove()
		e.DOM.Find("style").Remove()

		body := e.ChildText(".container .detail__body-text")
		imageURL := e.ChildAttr(".content__bg [dtr-evt=\"cover image\"] img", "src")
		image, err := u.GetImage(imageURL)
		if err != nil {
			fmt.Errorf("Can't fetch image: %v", err)
		}

		content = m.Content{
			Author:      e.ChildText(".detail__author"),
			FullTitle:   e.ChildText(".detail__title"),
			ImageURL:    imageURL,
			Image:       image,
			Content:     u.CleanText(body),
			PublishedAt: e.ChildText(".detail__date"),
			CrawledAt:   time.Now(),
		}

	})

	err := c.Visit(url)
	if err != nil {
		return content, err
	}

	return content, nil
}
