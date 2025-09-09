package scrapy

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/Adedunmol/scrapy/boards"
	"github.com/Adedunmol/scrapy/core"
	"github.com/google/uuid"
	"log"
	"sync"
)

func Collate(results <-chan []*core.Job) []*core.Job {
	// fetch the slices of jobs from the output channel and merge everything and return the result
	var scrapedJobs []*core.Job

	for jobs := range results {
		scrapedJobs = append(scrapedJobs, jobs...)
	}

	return scrapedJobs
}

func Coordinator(ctx context.Context, scheduled bool, searchTerm, location string, categoryStore categories.Store, jobStore jobs.JobStore) []*core.Job {
	fmt.Println("coordinator started")

	// fetch jobs from the categories table
	// 1. Fetch categories from DB
	categoriesData, err := categoryStore.GetCategories(ctx) // []string
	if err != nil {
		log.Printf("failed to fetch categories: %v", err)
		return nil
	}

	categoriesMap := make(map[string]uuid.UUID)

	for _, category := range categoriesData {
		categoriesMap[category.Name] = category.ID
	}

	//var wg sync.WaitGroup
	//
	//scrapers := []core.JobScraper{
	//	&boards.LinkedIn{
	//		BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
	//		JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
	//		Params: []struct{ Key, Value string }{
	//			{"location", location},
	//			{"keywords", searchTerm},
	//			{"f_TPR", "r86400"},
	//		},
	//	},
	//	&boards.JobberMan{
	//		BaseUrl: "https://www.jobberman.com/jobs",
	//		Params: []struct{ Key, Value string }{
	//			{"q", searchTerm},
	//		},
	//	},
	//}
	//
	//results := make(chan []*core.Job, 10)
	//
	//for _, scraper := range scrapers {
	//	wg.Add(1)
	//	go scraper.Run(&wg, results)
	//}

	var wg sync.WaitGroup
	results := make(chan []*core.Job, 50) // buffer for concurrent workers

	// Iterate over categories
	for _, category := range categoriesData {
		searchTerm := category.Name

		// Define scrapers for this search term
		scrapers := []core.JobScraper{
			&boards.LinkedIn{
				BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
				JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
				Params: []struct{ Key, Value string }{
					{"location", location},
					{"keywords", searchTerm},
					{"f_TPR", "r86400"},
				},
				Category:   category.Name,
				CategoryID: category.ID,
			},
			&boards.JobberMan{
				BaseUrl: "https://www.jobberman.com/jobs",
				Params: []struct{ Key, Value string }{
					{"q", searchTerm},
				},
				Category:   category.Name,
				CategoryID: category.ID,
			},
		}

		// Run scrapers concurrently
		for _, scraper := range scrapers {
			wg.Add(1)
			go func(scr core.JobScraper) {
				defer wg.Done()
				scr.Run(&wg, results)
			}(scraper)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	scrapedJobs := Collate(results)

	fmt.Println("scraped jobs: ", scrapedJobs)

	var bodies []jobs.CreateJobBody

	for _, job := range scrapedJobs {
		body := jobs.CreateJobBody{
			JobTitle:   job.Title,
			JobLink:    job.Link,
			DatePosted: job.DatePosted,
			CategoryID: job.CategoryID, // you pass this in when converting
			Origin:     "scraper",      // e.g. "LinkedIn" or "Jobberman"
			OriginID:   job.Id,         // use scraped Id as the origin ID
		}
		bodies = append(bodies, body)
	}

	jobStore.BatchCreateJobs(ctx, bodies)

	if scheduled {
		err := SendMail(core.Email, scrapedJobs)
		if err != nil {
			log.Printf("error occurred while sending mail: %v", errors.Unwrap(err))
		}
	}

	fmt.Println("coordinator finished")

	return scrapedJobs
}
