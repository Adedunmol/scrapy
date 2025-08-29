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

type JobberMan struct {
	BaseUrl string // https://www.jobberman.com/jobs/full-time?q=python&work_type=full-time
	JobUrl  string
	Params  []struct{ Key, Value string }
}

func (j *JobberMan) AddPagination(page int) string {
	dst := j.BuildUrl()

	return dst + fmt.Sprintf("&page=%d", page)
}

func (j *JobberMan) BuildUrl() string {
	fmt.Println("building url")

	params := make(url.Values)

	for _, v := range j.Params {
		if v.Value != "" {
			params.Add(v.Key, v.Value)
		}
	}

	if core.Location != "" {
		j.JobUrl = core.Location + "/" + j.JobUrl
	}

	dst := j.BaseUrl + "?" + params.Encode()

	fmt.Println("built url")

	return dst
}

func (j *JobberMan) ScrapeJobs(ctx context.Context, url string) ([]*core.Job, error) {
	fmt.Println("scraping url:", url)
	c := colly.NewCollector(
		colly.UserAgent(core.UserAgent),
	)

	var res []*core.Job
	var scrapeErr error

	// The main job card container
	c.OnHTML("div[data-cy='listing-cards-components']", func(e *colly.HTMLElement) {
		job := j.ParseJob(e)

		if &job != nil {
			res = append(res, &job)
		}
	})

	//c.OnError(func(r *colly.Response, err error) {
	//	fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	//	if r.StatusCode == 404 {
	//		// stop feeding the channel. close the channel.
	//		scrapeErr = scrapy.ErrNotFound
	//	}
	//	scrapeErr = err
	//})

	c.Visit(url)

	return res, scrapeErr
}

func (j *JobberMan) ParseJob(e *colly.HTMLElement) core.Job {
	var job core.Job

	job.Title = strings.TrimSpace(e.ChildText("a[data-cy='listing-title-link'] p"))
	job.Company = strings.TrimSpace(e.ChildText("p.text-sm.text-link-500 a"))
	job.Link = strings.TrimSpace(e.ChildAttr("a[data-cy='listing-title-link']", "href"))
	job.DatePosted = strings.TrimSpace(e.ChildText("div.flex.flex-row.items-center p.ml-auto"))

	spans := []string{}
	e.ForEach("div.flex.flex-wrap.mt-3.text-sm.text-gray-500.md\\:py-0 span", func(_ int, span *colly.HTMLElement) {
		spans = append(spans, strings.TrimSpace(span.Text))
	})

	if len(spans) >= 1 {
		job.Location = spans[0]
	}

	return job
}

func (j *JobberMan) FetchJobDetails(jobID string) (string, string) {

	return "", ""
}

func (j *JobberMan) Run(globalWg *sync.WaitGroup, results chan<- []*core.Job) {
	defer globalWg.Done()
	pagesCh := make(chan int, core.Buffer)
	var wg sync.WaitGroup

	// Spin up workers
	for i := 0; i < core.Workers; i++ {
		wg.Add(1)
		go core.Worker(pagesCh, j, results, &wg)
	}

	// Feed pages
	for i := 1; i <= core.Pages; i++ {
		pagesCh <- i
	}
	close(pagesCh)

	// Wait for workers to finish
	wg.Wait()
}
