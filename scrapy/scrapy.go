package scrapy

import (
	"context"
	"github.com/gocolly/colly"
	"sync"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

type Job struct {
	Id         string
	Title      string
	Company    string
	DatePosted string
	Location   string
	Applicants string
	Link       string
}

type Entry struct {
	Params []struct{ Key, Value string }
}

type JobScraper interface {
	BuildUrl() string
	ScrapeJobs(ctx context.Context, url string) ([]*Job, error)
	ParseJob(e *colly.HTMLElement) Job
	FetchJobDetails(jobID string) (string, string)
	Run(wg *sync.WaitGroup, results chan<- []*Job)
}

func Worker(pagesCh <-chan int, scraper JobScraper, results chan<- []*Job, wg *sync.WaitGroup) {
	// get the page from the channel and scrape. pass the pointer to the slice back to the output channel
	//ctx := context.Background()
	defer wg.Done()

	//for page := range pagesCh {

	//entry.Params = append(entry.Params, struct{ Key, Value string }{
	//	"start", strconv.Itoa(page),
	//})
	//dst := scraper.BuildUrl(entry)
	//
	//scrapedJobs, err := scraper.ScrapeJobs(ctx, dst)
	//if err != nil {
	//	log.Printf("scrape failed: %v\n", errors.Unwrap(err))
	//	return
	//}

	//results <- scrapedJobs
	//}
}

func Collate(results <-chan []*Job) []*Job {
	// fetch the slices of jobs from the output channel and merge everything and return the result
	var scrapedJobs []*Job

	for jobs := range results {
		scrapedJobs = append(scrapedJobs, jobs...)
	}

	return scrapedJobs
}
