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

type JobberMan struct {
	BaseUrl string // https://www.jobberman.com/jobs/full-time?q=python&work_type=full-time
	JobUrl  string
}

func (j *JobberMan) BuildUrl(entry *scrapy.Entry) string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range entry.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	dst := j.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (j *JobberMan) ScrapeJobs(ctx context.Context, url string) ([]*scrapy.Job, error) {
	fmt.Println("scraping url:", url)
	c := colly.NewCollector(
		colly.UserAgent(scrapy.UserAgent),
	)

	var res []*scrapy.Job

	// The main job card container
	c.OnHTML("div[data-cy='listing-cards-components']", func(e *colly.HTMLElement) {
		job := j.ParseJob(e)
		res = append(res, &job)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (j *JobberMan) ParseJob(e *colly.HTMLElement) scrapy.Job {
	var job scrapy.Job

	// Extract job title
	job.Title = strings.TrimSpace(e.ChildText("a[data-cy='listing-title-link'] p"))

	// Extract company name
	job.Company = strings.TrimSpace(e.ChildText("p.text-sm.text-link-500"))

	// Extract location (inside span tags with text-gray-700)
	job.Location = strings.TrimSpace(e.ChildText("div.flex.flex-wrap.mt-3 span:first-child"))

	// Extract link
	relativeLink := e.ChildAttr("a[data-cy='listing-title-link']", "href")
	job.Link = relativeLink // Already absolute on Jobberman

	return job
}

func (j *JobberMan) FetchJobDetails(jobID string) (string, string) {
	url := fmt.Sprintf("https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/%s", jobID)

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

const JobberManBuffer = 3
const JobberManWorkers = 3

func (j *JobberMan) Run(wg *sync.WaitGroup, results chan<- []*scrapy.Job) {
	pagesCh := make(chan int, JobberManBuffer)

	wg.Add(JobberManWorkers)
	for i := 0; i < JobberManWorkers; i++ {
		go scrapy.Worker(pagesCh, j, results, wg)
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
