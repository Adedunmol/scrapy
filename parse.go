package main

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/scrapy"
	"log"
	"strconv"
	"sync"
)

func worker(pagesCh <-chan int, results chan<- []*scrapy.Job, wg *sync.WaitGroup) {
	// get the page from the channel and scrape. pass the pointer to the slice back to the output channel
	ctx := context.Background()
	defer wg.Done()

	entry := &scrapy.Entry{Params: []struct{ Key, Value string }{
		{"location", Location},
		{"keywords", SearchTerm},
		//{"start", strconv.Itoa(page)},
		{SortKey, SortLast24Hours},
	}}

	for page := range pagesCh {
		entry.Params = append(entry.Params, struct{ Key, Value string }{
			"start", strconv.Itoa(page),
		})
		dst := BuildUrl(entry)

		scrapedJobs, err := ScrapeJobs(ctx, dst)
		if err != nil {
			log.Printf("scrape failed: %v\n", errors.Unwrap(err))
			return
		}

		results <- scrapedJobs
	}
}

func collate(results <-chan []*scrapy.Job) []*scrapy.Job {
	// fetch the slices of jobs from the output channel and merge everything and return the result
	var scrapedJobs []*scrapy.Job

	for jobs := range results {
		scrapedJobs = append(scrapedJobs, jobs...)
	}

	return scrapedJobs
}
