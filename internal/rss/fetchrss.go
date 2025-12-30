package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"time"
)

func verifyURL(feedURL string) (err error) {
	u, err := url.Parse(feedURL)

	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	return nil
}

func cleanFeed(rss *RSSFeed) {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, _ := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	err := verifyURL(feedURL)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unsafe url: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to create a request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{
		Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to process the request: %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading the response body: %w", err)
	}
	var rssfeed RSSFeed
	err = xml.Unmarshal(data, &rssfeed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error unmarshalling xml from rss feed: %w", err)
	}
	cleanFeed(&rssfeed)
	return &rssfeed, nil
}
