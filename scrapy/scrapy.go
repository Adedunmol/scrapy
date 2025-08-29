package scrapy

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/boards"
	"github.com/Adedunmol/scrapy/core"
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

func Coordinator(ctx context.Context, scheduled bool, searchTerm, location string) []*core.Job {
	fmt.Println("coordinator started")

	var wg sync.WaitGroup

	scrapers := []core.JobScraper{
		//&boards.GlassDoor{},
		&boards.LinkedIn{
			BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
			JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
			Params: []struct{ Key, Value string }{
				{"location", location},
				{"keywords", searchTerm},
				{"f_TPR", "r86400"},
			},
		},
		//&boards.Indeed{
		//	BaseUrl: "https://www.indeed.com/jobs",
		//	Params: []struct{ Key, Value string }{
		//		{"q", SearchTerm},
		//		{"l", Location},
		//		{"sort", "date"},
		//	},
		//},
		&boards.JobberMan{
			BaseUrl: "https://www.jobberman.com/jobs",
			Params: []struct{ Key, Value string }{
				{"q", searchTerm},
			},
		},
	}

	results := make(chan []*core.Job, 10)

	for _, scraper := range scrapers {
		wg.Add(1)
		go scraper.Run(&wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	scrapedJobs := Collate(results)

	fmt.Println("scraped jobs: ", scrapedJobs)

	if scheduled {
		err := SendMail(core.Email, scrapedJobs)
		if err != nil {
			log.Printf("error occurred while sending mail: %v", errors.Unwrap(err))
		}
	}

	fmt.Println("coordinator finished")

	return scrapedJobs
}
