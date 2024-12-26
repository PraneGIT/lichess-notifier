package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PraneGIT/lichess-notifier/internal/config"
	"github.com/PraneGIT/lichess-notifier/internal/fetcher"
	"github.com/PraneGIT/lichess-notifier/internal/models"
	"github.com/PraneGIT/lichess-notifier/internal/notifier"
	"github.com/PraneGIT/lichess-notifier/internal/scheduler"
)

func main() {
	log.Println("Starting the Lichess Notifier...")

	// env configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	emailNotifier := notifier.NewEmailNotifier(cfg.Email)

	results := make(chan models.Game, 10)
	errors := make(chan error, 10)
	done := make(chan struct{}) // Channel to signal goroutines to stop

	// Setup signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start fetchers -> each user
	for _, username := range cfg.Users {
		log.Printf("Initializing fetcher for user: %s\n", username)
		fetcher := fetcher.NewFetcher(cfg.LichessAPIBase, username)
		go func() {
			scheduler.Start(fetcher, results, errors)
			select {
			case <-done:
				return
			}
		}()
	}

	// Start email notification handler
	go func() {
		for {
			select {
			case game := <-results:
				log.Printf("New game fetched: %s vs %s\n", game.White, game.Black)
				emailBody := formatGameDetails(game)
				if err := emailNotifier.SendEmail("Game Lost Notification", emailBody); err != nil {
					log.Printf("Error sending email for game %s: %v\n", game.Site, err)
					errors <- err
				} else {
					log.Printf("Email sent for game %s\n", game.Site)
				}
			case <-done:
				return
			}
		}
	}()

	// Start error handler
	go func() {
		for {
			select {
			case err := <-errors:
				log.Printf("Error: %v\n", err)
			case <-done:
				return
			}
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("\nReceived signal: %v\n", sig)
	log.Println("Shutting down gracefully...")

	// Signal all goroutines to stop
	close(done)

	// Close channels
	close(results)
	close(errors)

	log.Println("Shutdown complete")
	os.Exit(0)
}

func formatGameDetails(game models.Game) string {
	return fmt.Sprintf(
		"Are you kidding ??? what the **** are you talking about MAN ? you arE a biggest looser i ever seen in my life ! you was doing pipI in your pamperS when i was beating players much more stronger then you!\n"+
			"Game Lost Notification:\n\n"+
			"Event: %s\n"+
			"Date: %s\n"+
			"White: %s\n"+
			"Black: %s\n"+
			"Result: %s\n"+
			"Termination: %s\n"+
			"\nGame Link: %s\n",
		game.Event,
		game.Date,
		game.White,
		game.Black,
		game.Result,
		game.Termination,
		game.Site,
	)
}
