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
	"time"
)

const Email = "oyewaleadedunmola@gmail.com"
const SearchTerm = "Python"
const Location = ""
const Workers = 3
const BufferSize = 3
const Pages = 10

func coordinator(ctx context.Context) {
	fmt.Println("coordinator started")

	scrapers := []scrapy.JobScraper{
		&boards.GlassDoor{},
		&boards.LinkedIn{
			BaseUrl: "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search",
			JobUrl:  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
			Params:  []struct{ Key, Value string }{
				//{"location", Location},
				//{"keywords", SearchTerm},
				//{"start", strconv.Itoa(page)},
				//{"f_TPR", "r86400"},
			},
		},
		&boards.Indeed{},
		&boards.JobberMan{},
	}

	results := make(chan []*scrapy.Job, len(scrapers))

	//var wg sync.WaitGroup
	//
	//pagesCh := make(chan int, BufferSize)
	//results := make(chan []*Job, BufferSize)
	//
	//wg.Add(Workers)
	//for i := 0; i < Workers; i++ {
	//	go worker(pagesCh, results, &wg)
	//}
	//
	//go func() {
	//	wg.Wait()
	//	close(results)
	//}()
	//
	//for i := 1; i <= Pages; i++ {
	//	pagesCh <- i
	//}
	//close(pagesCh)
	//

	_ = scrapy.Collate(results) // scrapedJobs

	//url := BuildUrl(entry)
	//scrapedJobs, err := ScrapeJobs(ctx, url)
	//if err != nil {
	//	log.Printf("scrape failed: %v\n", errors.Unwrap(err))
	//	return
	//}
	//
	//err = SendMail(Email, scrapedJobs)
	//if err != nil {
	//	log.Printf("error occurred while sending mail: %v", errors.Unwrap(err))
	//}
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
			1*time.Minute,
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
