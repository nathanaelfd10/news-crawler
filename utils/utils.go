package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	m "news-crawler/models"
)

func GetImage(url string) (m.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return m.Image{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return m.Image{}, fmt.Errorf("failed to fetch image: %s", response.Status)
	}

	result, err := io.ReadAll(response.Body)
	image := m.Image{
		URL:  url,
		Data: result,
	}
	if err != nil {
		return m.Image{}, err
	}

	return image, nil
}

func GetImages(urls []string) []m.Image {
	var images []m.Image
	for _, url := range urls {
		if image, err := GetImage(url); err == nil {
			images = append(images, image)
		} else {
			log.Println("can't get image:", err)
		}
	}
	return images
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	return text
}
