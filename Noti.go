package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// Constants
const (
	blogID         = "1359530524392796723"         // Replace with your Blogger blog ID
	apiKey         = "AIzaSyDfd1X3EloZnjY-I3COjIhSA3PeOwFJttQ"          // Replace with your Google API key
	telegramBotKey = "7940347256:AAFKJ8InC4rFJRx0InUyhub-dxZ9jfUNVJM"      // Replace with your Telegram bot token
	chatID         = "@Titu_Updates"         // Replace with your Telegram channel username or chat ID
	bloggerAPI     = "https://www.googleapis.com/blogger/v3/blogs/%s/posts"
	telegramAPI    = "https://api.telegram.org/bot%s/sendMessage"
	pollInterval   = 60 * time.Second // Time interval to check for new posts
)

// BlogPost represents a blog post from the Blogger API
type BlogPost struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// APIResponse represents the Blogger API response
type APIResponse struct {
	Items []BlogPost `json:"items"`
}

// fetchLatestPost fetches the latest blog post using the Blogger API
func fetchLatestPost() (*BlogPost, error) {
	client := resty.New()
	url := fmt.Sprintf(bloggerAPI, blogID)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"key": apiKey,
		}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error fetching blog posts: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(resp.Body(), &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	if len(apiResponse.Items) == 0 {
		return nil, nil // No posts available
	}

	return &apiResponse.Items[0], nil // Return the latest post
}

// sendToTelegram sends a blog post to the Telegram channel
func sendToTelegram(post *BlogPost) error {
	client := resty.New()
	url := fmt.Sprintf(telegramAPI, telegramBotKey)

	message := fmt.Sprintf("📢 *New Blog Published!*\n\n*Title:* %s\n[Read Here](%s)", post.Title, post.URL)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"chat_id":    chatID,
			"text":       message,
			"parse_mode": "Markdown",
		}).
		Post(url)

	if err != nil {
		return fmt.Errorf("error sending message to Telegram: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code from Telegram: %d", resp.StatusCode())
	}

	return nil
}

func main() {
	var lastPostURL string

	for {
		// Fetch the latest post
		post, err := fetchLatestPost()
		if err != nil {
			log.Printf("Error fetching latest post: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		// If there is a new post, send it to Telegram
		if post != nil && post.URL != lastPostURL {
			log.Printf("New blog detected: %s", post.Title)

			if err := sendToTelegram(post); err != nil {
				log.Printf("Error sending to Telegram: %v", err)
			} else {
				lastPostURL = post.URL // Update the last sent post URL
			}
		}

		time.Sleep(pollInterval) // Wait before polling again
	}
}
