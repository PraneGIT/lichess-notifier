package main

import (
	"github.com/PraneGIT/lichess-notifier/internal/config"
	"github.com/PraneGIT/lichess-notifier/internal/fetcher"

	"fmt"
)

func main() {
    fmt.Println("hello world!")

    cfg := config.LoadConfig()
    fetcher := fetcher.NewFetcher(cfg.LichessAPIBase, cfg.Users[0])

    games, err := fetcher.FetchGames()
    if err != nil {
        fmt.Println(err)
        return
    }

    if len(games) == 0 {
        fmt.Println("No games fetched")
        return
    }

	for game := range games {
		fmt.Println("----new game-----")
		fmt.Println(games[game].Moves)
	}
}

