package boards

import (
	"context"
	"fmt"
	"github.com/Adedunmol/scrapy/scrapy"
	"github.com/gocolly/colly"
	"net/url"
	"strings"
	"sync"
)

type GlassDoor struct {
	BaseUrl string // https://www.glassdoor.com/Job/jobs.htm?sc.keyword=golang&locT=C&locId=1147401&locKeyword=San%20Francisco,%20CA&p=2
	JobUrl  string
	Params  []struct{ Key, Value string }
}

func (g *GlassDoor) AddPagination(page int) string {
	dst := g.BuildUrl()

	return dst + fmt.Sprintf("&p=%d", page)
}

func (g *GlassDoor) BuildUrl() string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range g.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	dst := g.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (g *GlassDoor) ScrapeJobs(ctx context.Context, url string) ([]*scrapy.Job, error) {
	fmt.Println("scraping url:", url)
	c := colly.NewCollector(
		colly.UserAgent(scrapy.UserAgent),
	)

	var res []*scrapy.Job

	c.OnHTML("li.JobList__jobItem", func(e *colly.HTMLElement) {
		job := g.ParseJob(e)
		res = append(res, &job)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *GlassDoor) ParseJob(e *colly.HTMLElement) scrapy.Job {
	var job scrapy.Job

	job.Title = strings.TrimSpace(e.ChildText("a.JobCard_jobTitle"))
	job.Company = strings.TrimSpace(e.ChildText("div.JobCard_jobEmployerName"))
	job.Location = strings.TrimSpace(e.ChildText("div.JobCard_jobLocation"))

	relativeLink := e.ChildAttr("a.JobCard_jobTitle", "href")
	job.Link = "https://www.glassdoor.com" + relativeLink

	return job
}

func (g *GlassDoor) FetchJobDetails(jobID string) (string, string) {
	return "", ""
}

const GlassDoorBuffer = 3
const GlassDoorWorkers = 3

func (g *GlassDoor) Run(globalWg *sync.WaitGroup, results chan<- []*scrapy.Job) {
	defer globalWg.Done()
	pagesCh := make(chan int, LinkedInBuffer)
	var wg sync.WaitGroup

	// Spin up workers
	for i := 0; i < LinkedInWorkers; i++ {
		wg.Add(1)
		go scrapy.Worker(pagesCh, g, results, &wg)
	}

	// Feed pages
	for i := 1; i <= 10; i++ {
		pagesCh <- i
	}
	close(pagesCh)

	// Wait for workers to finish
	wg.Wait()
}
