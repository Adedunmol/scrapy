package boards

import (
	"context"
	"fmt"
	"github.com/Adedunmol/scrapy/core"
	"github.com/gocolly/colly"
	"net/url"
	"strings"
	"sync"
)

type Indeed struct {
	BaseUrl string // ?q=python&l=Texas
	JobUrl  string
	Params  []struct{ Key, Value string }
}

func (i *Indeed) AddPagination(page int) string {
	dst := i.BuildUrl()

	// pagination starts with nil, then 10 -> 20 -> 30
	page -= 1
	if page == 0 {
		return dst
	}

	page = page * 10

	return dst + fmt.Sprintf("&start=%d", page)
}

func (i *Indeed) BuildUrl() string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range i.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	dst := i.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (i *Indeed) ScrapeJobs(ctx context.Context, url string) ([]*core.Job, error) {
	fmt.Println("scraping url:", url)
	c := colly.NewCollector()
	c.UserAgent = core.UserAgent

	var res []*core.Job

	c.OnHTML("div.job_seen_beacon", func(e *colly.HTMLElement) {
		job := i.ParseJob(e)
		res = append(res, &job)
	})

	c.Visit(url)
	return res, nil
}

func (i *Indeed) ParseJob(e *colly.HTMLElement) core.Job {
	var job core.Job
	job.Title = strings.TrimSpace(e.ChildText("h2.jobTitle span"))
	job.Company = strings.TrimSpace(e.ChildText("span.companyName"))
	job.Location = strings.TrimSpace(e.ChildText("div.companyLocation"))
	job.DatePosted = strings.TrimSpace(e.ChildText("span.date"))

	relativeLink := e.ChildAttr("a.jcs-JobTitle", "href")
	job.Link = "https://www.indeed.com" + relativeLink

	return job
}

func (i *Indeed) FetchJobDetails(jobID string) (string, string) {
	return "", ""
}

func (i *Indeed) Run(globalWg *sync.WaitGroup, results chan<- []*core.Job) {
	defer globalWg.Done()
	pagesCh := make(chan int, core.Buffer)
	var wg sync.WaitGroup

	// Spin up workers
	for curr := 0; curr < core.Workers; curr++ {
		wg.Add(1)
		go core.Worker(pagesCh, i, results, &wg)
	}

	// Feed pages
	for i := 1; i <= 10; i++ {
		pagesCh <- i
	}
	close(pagesCh)

	// Wait for workers to finish
	wg.Wait()
}
