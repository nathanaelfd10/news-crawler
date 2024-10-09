package main

import (
	"context"
	"fmt"
	"log"
	cfg "news-crawler/config"
	db "news-crawler/database"
	m "news-crawler/models"
	u "news-crawler/utils"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Establish connection to MongoDB
	client, collection, err := db.CreateMongoDBConnection(cfg.MongoDBHost, cfg.MongoDBPort, cfg.DatabaseName, cfg.CollectionNameDetik)
	if err != nil {
		log.Fatalf("Error establishing MongoDB connection: %v", err)
	}
	defer client.Disconnect(context.TODO())

	c := colly.NewCollector()

	maxPages := cfg.DetikMaxPage // Set to 0 for infinite crawling
	currentPage := 0

	var articles []m.DetikArticle

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawling", r.URL.String())
	})

	c.OnHTML("article.list-content__item", func(e *colly.HTMLElement) {
		article := extractDetikArticle(e)
		articleURL := article.Article.URL
		fmt.Println("Found: ", article.Article.Preview.Title)

		content, err := GetContent(articleURL)
		if err != nil {
			fmt.Printf("Can't fetch article: %s. %v\n", articleURL, err)
		}

		article.Article.Content = content
		articles = append(articles, article)
	})

	for {
		currentPage++
		articles = nil // Reset articles at the start of each loop

		url := fmt.Sprintf("https://news.detik.com/indeks/%d", currentPage)

		fmt.Printf("Crawling page %d...\n", currentPage)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting page %d: %s\n", currentPage, err)
			break
		}

		if len(articles) == 0 {
			fmt.Println("No more article found in page.")
			break
		}

		// Save each articles to database.
		for _, article := range articles {
			_, err = collection.InsertOne(context.TODO(), article)
			if err != nil {
				log.Printf("\tFailed to insert article: %s\n", err)
				return
			}
		}

		if currentPage == maxPages {
			fmt.Printf("Reached max page %d\n", cfg.DetikMaxPage)
			break
		}
	}

	fmt.Println("Crawling completed.")
}

func extractDetikArticle(e *colly.HTMLElement) m.DetikArticle {
	thumbnailURL := e.ChildAttr(".media__image img", "src")
	thumbnail, err := u.GetImage(thumbnailURL)
	if err != nil {
		log.Printf("Failed to fetch image: %s\n", err)
	}

	article := m.DetikArticle{
		Article: m.Article[m.DetikExtraPreviewInfo, m.DetikExtraContentInfo]{
			URL: e.ChildAttr(".media__title > a", "href"),
			Preview: m.Preview[m.DetikExtraPreviewInfo]{
				Title:        e.ChildText(".media__title > a"),
				Thumbnail:    thumbnail.Data,
				ThumbnailURL: thumbnailURL,
				ExtraPreviewInfo: m.DetikExtraPreviewInfo{
					TimestampUTC: e.ChildAttr(".media__date > span", "d-time"),
					DateWIB:      e.ChildAttr(".media__date > span", "title"),
				},
				CrawledAt: time.Now(),
			},
		},
	}

	return article
}

func GetContent(url string) (m.Content[m.DetikExtraContentInfo], error) {
	var content m.Content[m.DetikExtraContentInfo]

	c := colly.NewCollector()

	c.OnHTML(".content__bg", func(e *colly.HTMLElement) {
		e.DOM.Find(".staticdetail_container, .lihatjg, .para_caption, .staticdetail_ads").Remove()
		e.DOM.Find("style").Remove()

		body := e.ChildText(".container .detail__body-text")

		// Fetch images
		var images []m.Image
		e.ForEach("figure.detail__media-image", func(_ int, el *colly.HTMLElement) {
			url := el.ChildAttr("img", "src")
			if url == "" {
				url = el.ChildAttr("img", "data-lazy")
			}

			image, err := u.GetImage(url)
			if err != nil {
				errMsg := fmt.Errorf("failed to get image from %s: %w", url, err)
				fmt.Println(errMsg)
			}

			if strings.Contains(url, "foto-news") {
				image.Caption = el.ChildText("figcaption p")
			} else {
				image.Caption = el.ChildText("figcaption")
			}
			images = append(images, image)
		})

		content = m.Content[m.DetikExtraContentInfo]{
			Author:      e.ChildText(".detail__author"),
			FullTitle:   e.ChildText(".detail__title"),
			Images:      images,
			Content:     u.CleanText(body),
			PublishedAt: e.ChildText(".detail__date"),
			CrawledAt:   time.Now(),
			ExtraContentInfo: m.DetikExtraContentInfo{
				Tags: e.ChildTexts(".detail__body-tag .nav .nav__item"),
			},
		}

	})

	err := c.Visit(url)
	if err != nil {
		return content, err
	}

	return content, nil
}
