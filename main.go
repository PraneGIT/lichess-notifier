package main

import (
	"fmt"
	"log"

	"github.com/PraneGIT/lichess-notifier/internal/config"
	"github.com/PraneGIT/lichess-notifier/internal/fetcher"
	"github.com/PraneGIT/lichess-notifier/internal/models"
	"github.com/PraneGIT/lichess-notifier/internal/notifier"
	"github.com/PraneGIT/lichess-notifier/internal/scheduler"
)

func main() {
	// Starting the program, log the event
	log.Println("Starting the Lichess Notifier...")

	cfg := config.LoadConfig()
	emailNotifier := notifier.NewEmailNotifier(cfg.Email)

	results := make(chan models.Game)
	errors := make(chan error)

	// Initialize fetchers and schedulers
	for _, username := range cfg.Users {
		log.Printf("Initializing fetcher for user: %s\n", username)
		fetcher := fetcher.NewFetcher(cfg.LichessAPIBase, username)
		go scheduler.Start(fetcher, results, errors)
	}

	// Goroutine to send email notifications
	go func() {
		for game := range results {
			log.Printf("New game fetched: %s vs %s\n", game.White, game.Black)
			emailBody := formatGameDetails(game)
			if err := emailNotifier.SendEmail("Game Lost Notification", emailBody); err != nil {
				log.Printf("Error sending email for game %s: %v\n", game.Site, err)
				errors <- err
			} else {
				log.Printf("Email sent for game %s\n", game.Site)
			}
		}
	}()

	// Goroutine to handle errors
	go func() {
		for err := range errors {
			log.Printf("Error: %v\n", err)
		}
	}()
	
	// Block indefinitely, keeping the program running
	select {}
}

func formatGameDetails(game models.Game) string {
	// Format the game details into a string
	return fmt.Sprintf(
		"Game Lost Notification:\n\n"+
			"Event: %s\n"+
			"Site: %s\n"+
			"Date: %s\n"+
			"White: %s\n"+
			"Black: %s\n"+
			"Result: %s\n"+
			"Termination: %s\n"+
			"\nGame Link: %s\n",
		game.Event,
		game.Site,
		game.Date,
		game.White,
		game.Black,
		game.Result,
		game.Termination,
		game.Site,
	)
}
