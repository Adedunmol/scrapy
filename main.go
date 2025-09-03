package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api"
	"github.com/Adedunmol/scrapy/database"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func runMigrations(db *pgxpool.Pool, dir string) error {
	goose.SetBaseFS(nil) // required if you're not embedding migrations

	sqlDB := stdlib.OpenDBFromPool(db)
	if sqlDB == nil {
		return fmt.Errorf("could not unwrap pgxpool to *sql.DB")
	}

	err := goose.Up(sqlDB, dir)
	if err != nil {
		// Check if it's a "no files found" or "dir not exist" error
		if errors.Is(err, os.ErrNotExist) {
			log.Println("No migration directory found, skipping...")
			return nil
		}
		if err.Error() == "no migration files found" {
			log.Println("No migration files found, skipping...")
			return nil
		}
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %s. relying on enviroment variables", err)
	}

	ctx := context.Background()

	pool, err := database.ConnectDB(ctx)

	defer pool.Close()
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	// Run migrations
	if err := runMigrations(pool, "./database/migrations"); err != nil {
		log.Fatal(fmt.Errorf("failed to run migrations: %w", err))
	}

	fmt.Println("Migrations applied successfully")

	// start job scheduler
	s, scheduleErr := gocron.NewScheduler()
	if scheduleErr != nil {
		// handle error
		err = fmt.Errorf("error creating new scheduler: %v", errors.Unwrap(scheduleErr))
		log.Fatal(err)
	}

	r := api.Routes(pool)

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
