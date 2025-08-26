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

type Indeed struct {
	BaseUrl string // https://www.indeed.com/jobs?q=python&l=Texas
	JobUrl  string
}

func (i *Indeed) BuildUrl(entry *scrapy.Entry) string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range entry.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	dst := i.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (i *Indeed) ScrapeJobs(ctx context.Context, url string) ([]*scrapy.Job, error) {
	fmt.Println("scraping url:", url)
	c := colly.NewCollector()
	c.UserAgent = scrapy.UserAgent

	var res []*scrapy.Job

	c.OnHTML("div.job_seen_beacon", func(e *colly.HTMLElement) {
		job := i.ParseJob(e)
		res = append(res, &job)
	})

	c.Visit(url)
	return res, nil
}

func (i *Indeed) ParseJob(e *colly.HTMLElement) scrapy.Job {
	var job scrapy.Job
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

const IndeedBuffer = 3
const IndeedWorkers = 3

func (i *Indeed) Run(wg *sync.WaitGroup, results chan<- []*scrapy.Job) {
	pagesCh := make(chan int, IndeedBuffer)

	wg.Add(JobberManWorkers)
	for curr := 0; curr < IndeedWorkers; curr++ {
		go scrapy.Worker(pagesCh, i, results, wg)
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
