package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/marekmchl/aggreGATOR/internal/database"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fetching failed - %v", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fetching failed - %v", err)
	}
	defer resp.Body.Close()

	rssData, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fetching failed - %v", err)
	}

	data := &RSSFeed{}
	if err := xml.Unmarshal(rssData, data); err != nil {
		return &RSSFeed{}, fmt.Errorf("fetching failed - %v", err)
	}

	// data.Channel.Title = html.UnescapeString(strings.TrimSpace(data.Channel.Title))
	// data.Channel.Description = html.UnescapeString(strings.TrimSpace(data.Channel.Description))

	return data, nil
}

func parseTime(timeString string) (time.Time, error) {
	pubTime, err := time.Parse(time.RFC1123, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC1123Z, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC3339, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC822, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC822Z, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC850, timeString)
	if err == nil {
		return pubTime, nil
	}
	pubTime, err = time.Parse(time.RFC3339Nano, timeString)
	if err == nil {
		return pubTime, nil
	}

	return time.Time{}, fmt.Errorf("unsupported time format")
}

func ScrapeFeeds(s *state.State) error {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("getting the feed was unsuccessful - %v", err)
	}
	feed, err = s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("marking the feed as fetched was unsuccessful - %v", err)
	}
	rssFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("fetching the feed was unsuccessful - %v", err)
	}
	for _, rssItem := range rssFeed.Channel.Item {
		pubTime, err := parseTime(rssItem.PubDate)
		if err != nil {
			return fmt.Errorf("parsing time was unsuccessful - %v", err)
		}
		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(strings.TrimSpace(rssItem.Title)),
			Url:         rssItem.Link,
			Description: html.UnescapeString(strings.TrimSpace(rssItem.Description)),
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		})
		if err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate key value violates unique constraint \"posts_url_key\"") {
			return fmt.Errorf("creating post was unsuccessful - %v", err)
		}
	}
	return nil
}
