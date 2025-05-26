package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

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

	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)

	return data, nil
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
		fmt.Printf("%v (%v, %v)\n", rssItem.Title, rssItem.PubDate, rssItem.Link)
		fmt.Println(rssItem.Description)
		fmt.Println()
	}
	return nil
}
