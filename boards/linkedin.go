package boards

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/scrapy"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"strings"
	"sync"
)

type LinkedIn struct {
	BaseUrl string // "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"
	JobUrl  string // https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/
}

func (l *LinkedIn) BuildUrl(entry *scrapy.Entry) string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range entry.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	dst := l.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (l *LinkedIn) ScrapeJobs(ctx context.Context, url string) ([]*scrapy.Job, error) {
	fmt.Println("scraping url:", url)
	// Instantiate default collector
	c := colly.NewCollector()

	c.UserAgent = scrapy.UserAgent

	var res []*scrapy.Job
	// On every a element which has href attribute call callback
	c.OnHTML("li div.base-card", func(e *colly.HTMLElement) {
		job := l.ParseJob(e)

		applicants, datePosted := l.FetchJobDetails(job.Id)

		job.Applicants = applicants
		job.DatePosted = datePosted

		res = append(res, &job)
	})

	c.Visit(url)

	return res, nil
}

func (l *LinkedIn) ParseJob(e *colly.HTMLElement) scrapy.Job {
	// parse jobs gotten from scraping
	var job scrapy.Job
	job.Id = strings.Split(e.Attr("data-entity-urn"), ":")[3]

	job.Title = strings.TrimSpace(e.ChildText("h3.base-search-card__title"))
	job.Company = strings.TrimSpace(e.ChildText("h4.base-search-card__subtitle"))
	job.Location = strings.TrimSpace(e.ChildText("span.job-search-card__location"))
	job.Link = e.ChildAttr("a.base-card__full-link", "href")

	return job
}

func (l *LinkedIn) FetchJobDetails(jobID string) (string, string) {
	url := l.JobUrl + jobID

	c := colly.NewCollector()

	c.UserAgent = scrapy.UserAgent

	var applicants, posted string

	c.OnHTML("span.num-applicants__caption.topcard__flavor--metadata.topcard__flavor--bullet", func(e *colly.HTMLElement) {
		applicants = strings.TrimSpace(e.Text)
	})

	c.OnHTML("span.posted-time-ago__text.topcard__flavor--metadata", func(e *colly.HTMLElement) {
		posted = strings.TrimSpace(e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		err = fmt.Errorf("error fetching job details (registered): %v", errors.Unwrap(err))
		log.Fatal(err)
	})

	c.Visit(url)

	return applicants, posted
}

const LinkedInBuffer = 3
const LinkedInWorkers = 3

func (l *LinkedIn) Run(wg *sync.WaitGroup, results chan<- []*scrapy.Job) {
	pagesCh := make(chan int, LinkedInBuffer)

	wg.Add(LinkedInWorkers)
	for i := 0; i < LinkedInWorkers; i++ {
		go scrapy.Worker(pagesCh, l, results, wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 1; i <= 10; i++ {
		pagesCh <- i
	}
	close(pagesCh)

}
