package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", response.Status)
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	return text
}
