package scrapy

import (
	"context"
	"github.com/gocolly/colly"
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
	BuildUrl(entry *Entry) string
	ScrapeJobs(ctx context.Context, url string) ([]*Job, error)
	ParseJob(e *colly.HTMLElement) Job
	FetchJobDetails(jobID string) (string, string)
}
