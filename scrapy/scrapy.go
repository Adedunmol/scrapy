package scrapy

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/boards"
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

func Collate(results <-chan []*Job) []*Job {
	// fetch the slices of jobs from the output channel and merge everything and return the result
	var scrapedJobs []*Job

	for jobs := range results {
		scrapedJobs = append(scrapedJobs, jobs...)
	}

	return scrapedJobs
}

func Coordinator(ctx context.Context, scheduled bool) {
	fmt.Println("coordinator started")

	var wg sync.WaitGroup

	scrapers := []JobScraper{
		//&boards.GlassDoor{},
		//&boards.LinkedIn{
		//	BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
		//	JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
		//	Params: []struct{ Key, Value string }{
		//		{"location", scrapy.Location},
		//		{"keywords", SearchTerm},
		//		{"f_TPR", "r86400"},
		//	},
		//},
		&boards.Indeed{
			BaseUrl: "https://www.indeed.com/jobs",
			Params: []struct{ Key, Value string }{
				{"q", SearchTerm},
				{"l", Location},
				{"sort", "date"},
			},
		},
		//&boards.JobberMan{
		//	BaseUrl: "https://www.jobberman.com/jobs",
		//	Params: []struct{ Key, Value string }{
		//		{"q", SearchTerm},
		//	},
		//},
	}

	results := make(chan []*Job, 10)

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
		err := SendMail(Email, scrapedJobs)
		if err != nil {
			log.Printf("error occurred while sending mail: %v", errors.Unwrap(err))
		}
	}

	fmt.Println("coordinator finished")
}
