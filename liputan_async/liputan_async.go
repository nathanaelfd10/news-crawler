package main

import (
	"context"
	"fmt"
	"log"
	cfg "news-crawler/config"
	db "news-crawler/database"
	m "news-crawler/models"
	u "news-crawler/utils"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func crawlIndex(collection *mongo.Collection) {
	log.Println("Crawling index..")

	c := colly.NewCollector()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	var articles []m.LiputanArticle

	c.OnHTML("[data-article-id]", func(e *colly.HTMLElement) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			article := m.LiputanArticle{
				Article: m.Article[m.LiputanExtraPreviewInfo, m.LiputanExtraContentInfo]{
					URL: fmt.Sprintf("https://www.liputan6.com/news/read/%s", e.Attr("data-article-id")),
					Preview: m.Preview[m.LiputanExtraPreviewInfo]{
						Title: e.Attr("data-title"),
					},
				},
			}

			fmt.Println("Found: ", article.Article.Preview.Title)

			content, err := GetContent(article.Article.URL)
			if err != nil {
				fmt.Printf("Can't fetch article: %s. %v\n", article.Article.URL, err)
			}

			article.Article.Content = content

			mutex.Lock()
			articles = append(articles, article)
			mutex.Unlock()
		}()
	})

	c.Visit("https://www.liputan6.com/")

	// Wait until all articles have finished crawling
	wg.Wait()

	if len(articles) == 0 {
		fmt.Println("No articles found in page.")
		return
	}

	saveArticlesToDb(articles, collection)

	log.Println("Done crawling index")
}

func crawlPagination(collection *mongo.Collection) {
	c := colly.NewCollector()

	maxPages := cfg.LiputanMaxPage // Set to 0 for infinite crawling
	currentPage := 0

	var articles []m.LiputanArticle

	var wg sync.WaitGroup
	var mutex sync.Mutex

	c.OnHTML("#indeks-articles article", func(e *colly.HTMLElement) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			article := extractLiputanArticle(e)
			fmt.Println("Found: ", article.Article.Preview.Title)

			content, err := GetContent(article.Article.URL)
			if err != nil {
				fmt.Printf("Can't fetch article: %s. %v\n", article.Article.URL, err)
			}

			article.Article.Content = content

			mutex.Lock()
			articles = append(articles, article)
			mutex.Unlock()
		}()
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

		wg.Wait()

		if len(articles) == 0 {
			fmt.Println("No more article found in page.")
			break
		}

		saveArticlesToDb(articles, collection)

		if currentPage == maxPages {
			fmt.Printf("Reached max page (%d)\n", cfg.LiputanMaxPage)
			break
		}
	}
}

func saveArticlesToDb(articles []m.LiputanArticle, collection *mongo.Collection) {
	for _, article := range articles {
		_, err := collection.InsertOne(context.TODO(), article)
		if err != nil {
			log.Printf("\tFailed to insert article: %s", err)
			return
		}
	}
}

func main() {
	client, collection, err := db.CreateMongoDBConnection(cfg.MongoDBHost, cfg.MongoDBPort, cfg.DatabaseName, cfg.CollectionNameLiputan6)
	if err != nil {
		log.Fatalf("Error establishing MongoDB connection: %v", err)
	}
	defer client.Disconnect(context.TODO())

	log.Println("Crawling index..")
	crawlIndex(collection)

	log.Println("Crawling pagination..")
	crawlPagination(collection)

	fmt.Println("Crawling finished.")
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
				Thumbnail:    thumbnail.Data,
				ThumbnailURL: thumbnail.URL,
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

		image.Caption = e.ChildText(".read-page--top-media figcaption.read-page--photo-gallery--item__caption")

		content = m.Content[m.LiputanExtraContentInfo]{
			Author:      e.ChildText("[class=\"read-page--header--author__name fn\"]"),
			FullTitle:   e.ChildText("h1[itemprop=\"headline\"]"),
			Images:      []m.Image{image},
			Content:     u.CleanText(body),
			PublishedAt: e.ChildText("time[itemprop=\"datePublished\"]"),
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
