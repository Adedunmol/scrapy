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
	dst := fmt.Sprintf("https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/%s", jobID)

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

	c.Visit(dst)

	return applicants, posted
}
