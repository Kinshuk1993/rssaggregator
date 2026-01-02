package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kinshuk1993/rssaggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)

	// use a new ticker
	ticker := time.NewTicker(timeBetweenRequests)

	// every time a new tick comes in, start scraping
	for ; ; <- ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Printf("Error fetching feeds to scrape: %v", err)
			// this function to fetch should never stop working/scraping, so continue
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(&wg, db, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Error marking feed %d as fetched: %v", feed.ID, err)
		return
	}

	// actually scrape the feed
	rssFeed, err := URLToFeed(feed.Url)
	if err != nil {
		log.Printf("Error scraping feed %d from url %s: %v", feed.ID, feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Printf("Scraped post %s on feed %s", item.Title, feed.Name)
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published date %s for post %s: %v", item.PubDate, item.Title, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url: item.Link,
			FeedID: feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Post %s for feed %s already exists, skipping", item.Title, feed.Name)
				continue
			}
			log.Printf("Error saving post %s for feed %s: %v", item.Title, feed.Name, err)
		}
	}

	log.Printf("Feed %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}