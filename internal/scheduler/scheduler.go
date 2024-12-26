package scheduler

import(
	"github.com/PraneGIT/lichess-notifier/internal/fetcher"
	"github.com/PraneGIT/lichess-notifier/internal/models"
	"time"
    "sync"
    "log"
)

//stores IDs of already processed games
type GameTracker struct {
    processed map[string]struct{}
    mu        sync.Mutex
}

func NewGameTracker() *GameTracker {
    return &GameTracker{processed: make(map[string]struct{})}
}

func (gt *GameTracker) IsNewGame(gameID string) bool {
    gt.mu.Lock()
    defer gt.mu.Unlock()

    if _, exists := gt.processed[gameID]; exists {
        return false
    }
    gt.processed[gameID] = struct{}{}
    return true
}

func Start(fetcher *fetcher.Fetcher, results chan<- models.Game, errors chan<- error) {
    tracker := NewGameTracker()
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    log.Println("Start function initiated. Fetcher started.")

    for range ticker.C {
        log.Println("Fetching games...")

        games, err := fetcher.FetchGames()
        if err != nil {
            log.Printf("Error fetching games: %v\n", err)
            errors <- err
            continue
        }

        log.Printf("Successfully fetched %d games\n", len(games))

        for _, game := range games {
            if tracker.IsNewGame(game.Site) && isLoss(game, fetcher.User) {
                log.Printf("New loss game found: %s\n", game.Site)
                results <- game
            }
        }
    }
}


func isLoss(game models.Game, user string) bool {
    return (game.White == user && game.Result == "0-1") || (game.Black == user && game.Result == "1-0")
}