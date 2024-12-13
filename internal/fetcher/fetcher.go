package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PraneGIT/lichess-notifier/internal/models"
)

type Fetcher struct {
    apiBase string
    user    string
}

func NewFetcher(apiBase, user string) *Fetcher {
    return &Fetcher{apiBase: apiBase, user: user}
}

func (f *Fetcher) FetchGames() ([]models.Game, error) {
    url := fmt.Sprintf("%s%s?max=2", f.apiBase, f.user)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return parseGames(string(body)), nil
}

func parseGames(data string) []models.Game {
    games := []models.Game{}
    gameBlocks := strings.Split(data, "\\n\\n")

    for _, block := range gameBlocks {
        if strings.TrimSpace(block) == "" {
            continue
        }

        game := models.Game{
            Event:       extractField(block, "Event"),
            Site:        extractField(block, "Site"),
            Date:        extractField(block, "Date"),
            White:       extractField(block, "White"),
            Black:       extractField(block, "Black"),
            Result:      extractField(block, "Result"),
            Termination: extractField(block, "Termination"),
        }
        games = append(games, game)
    }

    return games
}

func extractField(block, field string) string {
    regex := regexp.MustCompile(fmt.Sprintf(`\\[%s "(.*?)"\\]`, field))
    match := regex.FindStringSubmatch(block)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}