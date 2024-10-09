package models

import "time"

type ArticlePreview[T any] struct {
	Title           string    `bson:"title"`
	URL             string    `bson:"url"`
	Thumbnail       []byte    `bson:"thumbnail"`
	ThumbnailURL    string    `bson:"thumbnail_url"`
	WebSpecificInfo T         `bson:"web_specific_info"`
	CrawledAt       time.Time `bson:"crawled_at"`
}

type DetikSpecificInfo struct {
	TimestampUTC string `bson:"timestamp_utc"`
	DateWIB      string `bson:"date_wib"`
}

type LiputanSpecificInfo struct {
	UpdatedAt   string `bson:"updated_at"`
	Description string `bson:"description"`
	Category    string `bson:"category"`
}

type Content struct {
	FullTitle   string    `bson:"full_title"`
	Content     string    `bson:"content"`
	Author      string    `bson:"author"`
	ImageURL    string    `bson:"image_url"`
	Image       []byte    `bson:"image"`
	PublishedAt string    `bson:"published_at"`
	CrawledAt   time.Time `bson:"crawled_at"`
}

type ArticleDetik struct {
	ArticlePreview ArticlePreview[DetikSpecificInfo] `bson:"article_preview"`
	Content
}

// type BaseInfo struct {
// 	Title        string  `bson:"title"`
// 	URL          string  `bson:"url"`
// 	Thumbnail    []byte  `bson:"thumbnail"`
// 	ThumbnailURL string  `bson:"thumbnail_url"`
// 	Content      Content `bson:"content"`
// 	CrawledAt    string  `bson:"crawled_at"`
// }

// type ArticleDetik struct {
// 	BaseInfo
// 	Content
// 	TimestampUTC string `bson:"timestamp_utc"`
// 	DateWIB      string `bson:"date_wib"`
// }

// type ArticleLiputan struct {
// 	BaseInfo
// 	Content
// 	UpdatedAt   string `bson:"updated_at"`
// 	Description string `bson:"description"`
// 	Category    string `bson:"category"`
// }

// type ArticleLiputan struct {
// 	Title        string    `bson:"title"`
// 	URL          string    `bson:"url"`
// 	Thumbnail    []byte    `bson:"thumbnail"`
// 	ThumbnailURL string    `bson:"thumbnail_url"`
// 	Description  string    `bson:"description"`
// 	Category     string    `bson:"category"`
// 	Content      Content   `bson:"content"`
// 	PublishedAt  string    `bson:"published_at"`
// 	CrawledAt    time.Time `bson:"crawled_at"`
// }
