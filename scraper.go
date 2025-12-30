package main

import (
	"context"
	"log"
	"sync"
	"time"

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
	}

	log.Printf("Feed %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}