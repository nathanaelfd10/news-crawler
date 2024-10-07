package models

import (
	"time"
)

type Preview[P any] struct {
	Title            string    `bson:"title"`
	Thumbnail        []byte    `bson:"thumbnail"`
	ThumbnailURL     string    `bson:"thumbnail_url"`
	ExtraPreviewInfo P         `bson:"extra_preview_info"`
	CrawledAt        time.Time `bson:"crawled_at"`
}

type DetikExtraPreviewInfo struct {
	TimestampUTC string `bson:"timestamp_utc"`
	DateWIB      string `bson:"date_wib"`
}

type DetikExtraContentInfo struct {
	Tags []string `bson:"tags"`
}

type LiputanExtraPreviewInfo struct {
	Description string `bson:"description"`
	Category    string `bson:"category"`
}

type LiputanExtraContentInfo struct {
	UpdatedAt string `bson:"updated_at"`
}

type Content[T any] struct {
	FullTitle        string    `bson:"full_title"`
	Content          string    `bson:"content"`
	Author           string    `bson:"author"`
	ImageURL         string    `bson:"image_url"`
	Image            []byte    `bson:"image"`
	ImageCaption     string    `bson:"image_caption"`
	ExtraContentInfo T         `bson:"extra_content_info"`
	PublishedAt      string    `bson:"published_at"`
	CrawledAt        time.Time `bson:"crawled_at"`
}

type Article[P any, C any] struct {
	URL     string     `bson:"url"`
	Preview Preview[P] `bson:"preview"`
	Content Content[C] `bson:"content"`
}

type DetikArticle struct {
	Article Article[DetikExtraPreviewInfo, DetikExtraContentInfo] `bson:"article"`
}

type LiputanArticle struct {
	Article Article[LiputanExtraPreviewInfo, LiputanExtraContentInfo] `bson:"article"`
}
