package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	fmt.Println("entered")

	ctx := context.Background()
	// start job scheduler
	s, scheduleErr := gocron.NewScheduler()
	if scheduleErr != nil {
		// handle error
		err = fmt.Errorf("error creating new scheduler: %v", errors.Unwrap(scheduleErr))
		//log.Fatal(err)
		fmt.Println("nigga")
	}

	fmt.Println("entered2")

	r := api.Routes()
	//
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: r}

	go func() {
		log.Printf("starting web server on port %s", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(fmt.Errorf("error starting web server on port %s: %w", port, err))
		}
	}()

	// register the function to be executed (run coordinator)
	_, err = s.NewJob(
		gocron.DurationJob(
			1*time.Minute,
		),
		//gocron.NewTask(
		//	scrapy.Coordinator,
		//	ctx,
		//	true,
		//	scrapy.SearchTerm,
		//	scrapy.Location,
		//),
		gocron.NewTask(
			func(ctx context.Context) {
				fmt.Println("running")
			},
			ctx,
		),
	)
	if err != nil {
		err = fmt.Errorf("error adding job to scheduler: %v", errors.Unwrap(err))
		log.Fatal(err)
	}

	fmt.Println("entered3")
	// start the scheduler
	s.Start()

	// block forever until you receive a CTRL-C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	fmt.Println("entered4")
	<-stop
	fmt.Println("entered6")

	err = s.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("error shutting down scheduler: %v", errors.Unwrap(err)))
	}

	fmt.Println("entered5")
}
