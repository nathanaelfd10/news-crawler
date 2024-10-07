package main

import (
	"context"
	cfg "detik-scraper/config"
	db "detik-scraper/database"
	m "detik-scraper/models"
	u "detik-scraper/utils"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	client, collection, err := db.CreateMongoDBConnection(cfg.MongoDBHost, cfg.MongoDBPort, cfg.DatabaseName, cfg.CollectionNameLiputan6)
	if err != nil {
		log.Fatalf("Error establishing MongoDB connection: %v", err)
	}
	defer client.Disconnect(context.TODO())

	c := colly.NewCollector()

	maxPages := cfg.LiputanMaxPage // Set to 0 for infinite crawling
	currentPage := 0

	var articles []m.LiputanArticle

	c.OnHTML("#indeks-articles article", func(e *colly.HTMLElement) {
		article := extractLiputanArticle(e)
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
		articles = nil

		url := fmt.Sprintf("https://www.liputan6.com/news/indeks?page=%d", currentPage)

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
				log.Printf("\tFailed to insert article: %s", err)
				return
			}
		}

		if currentPage == maxPages {
			fmt.Printf("Reached max page %d\n", cfg.LiputanMaxPage)
			break
		}

	}

	fmt.Println("Crawling completed.")
}

func extractLiputanArticle(e *colly.HTMLElement) m.LiputanArticle {
	thumbnailURL := e.ChildAttr("[data-template-var=\"image\"]", "src")
	thumbnail, err := u.GetImage(thumbnailURL)
	if err != nil {
		log.Printf("Failed to fetch image: %s", err)
	}

	article := m.LiputanArticle{
		Article: m.Article[m.LiputanExtraPreviewInfo, m.LiputanExtraContentInfo]{
			URL: e.ChildAttr("[data-template-var=\"url\"]", "href"),
			Preview: m.Preview[m.LiputanExtraPreviewInfo]{
				Title:        e.ChildText("[data-template-var=\"title\"]"),
				Thumbnail:    thumbnail,
				ThumbnailURL: thumbnailURL,
				ExtraPreviewInfo: m.LiputanExtraPreviewInfo{
					Description: e.ChildText("[data-template-var=\"summary\"]"),
					Category:    e.ChildText("[data-template-var=\"category\"]"),
				},
				CrawledAt: time.Now(),
			},
		},
	}

	return article
}

func GetContent(url string) (m.Content[m.LiputanExtraContentInfo], error) {
	var content m.Content[m.LiputanExtraContentInfo]

	c := colly.NewCollector()

	c.OnHTML(".inner-container-article", func(e *colly.HTMLElement) {
		e.DOM.Find(".baca-juga-collections").Remove()
		e.DOM.Find(".baca-juga-collections__detail").Remove()

		body := e.ChildText(".inner-container-article .article-content-body")
		imageURL := e.ChildAttr(".read-page--top-media .read-page--photo-gallery--item__picture img", "src")
		image, err := u.GetImage(imageURL)
		if err != nil {
			log.Printf("Can't fetch image: %s\n", err)
		}

		content = m.Content[m.LiputanExtraContentInfo]{
			Author:       e.ChildText("[class=\"read-page--header--author__name fn\"]"),
			FullTitle:    e.ChildText("h1[itemprop=\"headline\"]"),
			ImageURL:     imageURL,
			Image:        image,
			ImageCaption: e.ChildText(".read-page--top-media figcaption[class=\"read-page--photo-gallery--item__caption\"]"),
			Content:      u.CleanText(body),
			PublishedAt:  e.ChildText("time[itemprop=\"datePublished\"]"),
			ExtraContentInfo: m.LiputanExtraContentInfo{
				UpdatedAt: e.ChildText("time[itemprop=\"dateModified\"]"),
			},
			CrawledAt: time.Now(),
		}
	})

	err := c.Visit(url)
	if err != nil {
		return content, err
	}

	return content, nil
}
