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
	mongoDBURI := fmt.Sprintf("mongodb://%s:%s", c.MongoDBHost, c.MongoDBPort)
	client, collection, err := db.ConnectToMongoDB(mongoDBURI, c.DatabaseName, c.CollectionNameLiputan6)
	if err != nil {
		log.Fatalf("Error establishing MongoDB connection: %v", err)
	}
	defer client.Disconnect(context.TODO())

	c := colly.NewCollector()

	maxPages := 0 // Set to 0 for infinite crawling
	currentPage := 0

	var articles []m.ArticleLiputan

	c.OnHTML("#indeks-articles article", func(e *colly.HTMLElement) {
		article := extractLiputanArticle(e)
		fmt.Println("Found: ", article.Title)

		content, err := GetContent(article.URL)
		if err != nil {
			fmt.Println("Can't fetch article: %s. %v", article.URL, err)
		}

		article.Content = content

		articles = append(articles, article)
	})

	for {
		currentPage++

		url := fmt.Sprintf("https://www.liputan6.com/news/indeks?page=%d", currentPage)

		fmt.Printf("Crawling page %d...\n", currentPage)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting page %d: %s\n", currentPage, err)
			break
		}

		// Finish crawling when these conditions are met.
		if len(articles) == 0 || currentPage == maxPages {
			fmt.Println("No article found in page.")
			break
		} else {
			// Save each articles to database.
			for _, article := range articles {
				_, err = collection.InsertOne(context.TODO(), article)
				if err != nil {
					log.Printf("\tFailed to insert article: %s", err)
					return
				}
			}

			// Reset articles.
			articles = nil
		}
	}

	fmt.Println("Crawling completed.")
}

func extractLiputanArticle(e *colly.HTMLElement) m.ArticleLiputan {
	thumbnailURL := e.ChildAttr("[data-template-var=\"image\"]", "src")
	thumbnail, err := u.GetImage(thumbnailURL)
	if err != nil {
		fmt.Errorf("Failed to fetch image: %s", err)
	}

	article := m.ArticleLiputan{
		Title:        e.ChildText("[data-template-var=\"title\"]"),
		URL:          e.ChildAttr("[data-template-var=\"url\"]", "href"),
		Description:  e.ChildText("[data-template-var=\"summary\"]"),
		ThumbnailURL: thumbnailURL,
		Thumbnail:    thumbnail,
		Category:     e.ChildText("[data-template-var=\"category\"]"),
		PublishedAt:  e.ChildAttr("[data-template-var=\"date\"] time", "datetime"),
		CrawledAt:    time.Now(),
	}

	return article
}

func GetContent(url string) (m.Content, error) {
	var content m.Content

	c := colly.NewCollector()

	c.OnHTML(".inner-container-article", func(e *colly.HTMLElement) {
		e.DOM.Find(".baca-juga-collections").Remove()
		e.DOM.Find(".baca-juga-collections__detail").Remove()

		body := e.ChildText(".inner-container-article .article-content-body")
		imageURL := e.ChildAttr(".read-page--top-media .read-page--photo-gallery--item__picture img", "src")
		image, err := u.GetImage(imageURL)
		if err != nil {
			fmt.Errorf("Can't fetch image: %v", err)
		}

		content = m.Content{
			Author:      e.ChildText("[class=\"read-page--header--author__name fn\"]"),
			FullTitle:   e.ChildText("h1[itemprop=\"headline\"]"),
			ImageURL:    imageURL,
			Image:       image,
			Content:     u.CleanText(body),
			PublishedAt: e.ChildText("time[itemprop=\"datePublished\"]"),
			UpdatedAt:   e.ChildText("time[itemprop=\"dateModified\"]"),
			CrawledAt:   time.Now(),
		}
	})

	err := c.Visit(url)
	if err != nil {
		return content, err
	}

	return content, nil
}
