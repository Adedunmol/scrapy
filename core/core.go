package core

import (
	"context"
	"errors"
	"github.com/gocolly/colly"
	"log"
	"sync"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
const Location = ""
const Workers = 3
const Buffer = 5
const Pages = 10
const Email = "oyewaleadedunmola@gmail.com"
const SearchTerm = "Python"

var ErrNotFound = errors.New("page (jobs) not found")

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
	AddPagination(page int) string
	ScrapeJobs(ctx context.Context, url string) ([]*Job, error)
	ParseJob(e *colly.HTMLElement) Job
	FetchJobDetails(jobID string) (string, string)
	Run(globalWg *sync.WaitGroup, results chan<- []*Job)
}

func Worker(pagesCh <-chan int, scraper JobScraper, results chan<- []*Job, wg *sync.WaitGroup) {
	// get the page from the channel and scrape. pass the pointer to the slice back to the output channel
	ctx := context.Background()
	defer wg.Done()

	for page := range pagesCh {

		dst := scraper.AddPagination(page)

		scrapedJobs, err := scraper.ScrapeJobs(ctx, dst)
		if err != nil {
			log.Printf("scrape failed: %v\n", errors.Unwrap(err))
			//if errors.Is(err, ErrNotFound) { // define ErrNotFound for 404
			//	// signal the controller to stop
			//	select {
			//	case stopCh <- struct{}{}:
			//	default:
			//	}
			//	return
			//}
			continue
			//return
		}

		results <- scrapedJobs
	}
}
