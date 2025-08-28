package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/boards"
	"github.com/Adedunmol/scrapy/scrapy"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

const Email = "oyewaleadedunmola@gmail.com"
const SearchTerm = "Python"
const Workers = 3
const BufferSize = 3

func coordinator(ctx context.Context) {
	fmt.Println("coordinator started")

	var wg sync.WaitGroup

	scrapers := []scrapy.JobScraper{
		//&boards.GlassDoor{},
		&boards.LinkedIn{
			BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
			JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
			Params: []struct{ Key, Value string }{
				{"location", scrapy.Location},
				{"keywords", SearchTerm},
				{"f_TPR", "r86400"},
			},
		},
		//&boards.Indeed{},
		&boards.JobberMan{
			BaseUrl: "https://www.jobberman.com/jobs",
			Params: []struct{ Key, Value string }{
				{"q", SearchTerm},
			},
		},
	}

	results := make(chan []*scrapy.Job, 10)

	for _, scraper := range scrapers {
		wg.Add(1)
		go scraper.Run(&wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	scrapedJobs := scrapy.Collate(results)

	fmt.Println("scraped jobs: ", scrapedJobs)

	err := SendMail(Email, scrapedJobs)
	if err != nil {
		log.Printf("error occurred while sending mail: %v", errors.Unwrap(err))
	}
	fmt.Println("coordinator finished")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	ctx := context.Background()
	// start job scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		err = fmt.Errorf("error creating new scheduler: %v", errors.Unwrap(err))
		log.Fatal(err)
	}

	// register the function to be executed (run coordinator)
	_, err = s.NewJob(
		gocron.DurationJob(
			3*time.Minute,
		),
		gocron.NewTask(
			coordinator,
			ctx,
		),
	)
	if err != nil {
		err = fmt.Errorf("error adding job to scheduler: %v", errors.Unwrap(err))
		log.Fatal(err)
	}

	// start the scheduler
	s.Start()

	// block forever until you receive a CTRL-C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	err = s.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("error shutting down scheduler: %v", errors.Unwrap(err)))
	}
}
