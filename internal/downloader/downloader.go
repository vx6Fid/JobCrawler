package downloader

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Downloader struct {
	collector *colly.Collector
}

func NewDownloader() *Downloader {
	c := colly.NewCollector(
		colly.Async(true),
		// colly.AllowedDomains("weworkremotely.com", "amazon.jobs", "linkedin.com"),
	)
	return &Downloader{collector: c}
}

func (d *Downloader) FetchWithParser(ctx context.Context, url string, parseFunc func(e *colly.HTMLElement)) error {
	done := make(chan struct{})
	var resultErr error

	// Register the parsing handler
	d.collector.OnHTML("body", parseFunc)

	// Mark crawl finished
	d.collector.OnScraped(func(_ *colly.Response) {
		select {
		case <-done:
			// already closed due to error
		default:
			close(done)
		}
	})

	// Start crawl in background
	go func() {
		resultErr = d.collector.Visit(url)
		d.collector.Wait()
		select {
		case <-done:
		default:
			close(done)
		}
	}()

	// Timeout or success
	select {
	case <-ctx.Done():
		log.Printf("--- [TIMEOUT] --- Fetch cancelled for: %s", url)
		return fmt.Errorf("fetch timeout: %w", ctx.Err())
	case <-done:
		return resultErr
	}
}
