package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
    "github.com/notnil/chess"

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
    fmt.Println(url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-ndjson")

    client := &http.Client{}
    resp, err := client.Do(req)
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
	var games []models.Game
	pgns := strings.Split(data, "\n\n\n") //games by double newlines

	for _, pgn := range pgns {
		if strings.TrimSpace(pgn) == "" {
			continue
		}

		// reader := strings.NewReader(pgn)
		game := chess.NewGame()
		if err := game.UnmarshalText([]byte(pgn)); err != nil {
			fmt.Println("Error parsing PGN:", err)
			continue
		}

		headers := extractHeaders(pgn)

		games = append(games, models.Game{
            Event:            headers["Event"],
            Site:             headers["Site"],
            Date:             headers["Date"],
            White:            headers["White"],
            Black:            headers["Black"],
            Result:           headers["Result"],
            UTCDate:          headers["UTCDate"],
            UTCTime:          headers["UTCTime"],
            WhiteElo:         headers["WhiteElo"],
            BlackElo:         headers["BlackElo"],
            WhiteRatingDiff:  headers["WhiteRatingDiff"],
            BlackRatingDiff:  headers["BlackRatingDiff"],
            Variant:          headers["Variant"],
            TimeControl:      headers["TimeControl"],
            ECO:              headers["ECO"],
            Termination:      headers["Termination"],
            Moves:            game.String(), // PGN moves as a string
        })
        
	}
	return games
}

func extractHeaders(pgn string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(pgn, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, " \"") {
			parts := strings.SplitN(line[1:len(line)-1], " \"", 2)
			if len(parts) == 2 {
				headers[parts[0]] = strings.Trim(parts[1], "\"")
			}
		}
	}
	return headers
}

func extractField(block, field string) string {
    regex := regexp.MustCompile(fmt.Sprintf(`\\[%s "(.*?)"\\]`, field))
    match := regex.FindStringSubmatch(block)
    if len(match) > 1 {
        return match[1]
    }
    return ""
}