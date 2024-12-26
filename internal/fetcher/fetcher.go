package fetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"github.com/notnil/chess"

	"github.com/PraneGIT/lichess-notifier/internal/models"
)

type Fetcher struct {
    apiBase string
    User    string
}

func NewFetcher(apiBase, user string) *Fetcher {
    return &Fetcher{apiBase: apiBase, User: user}
}

func (f *Fetcher) FetchGames() ([]models.Game, error) {
    url := fmt.Sprintf("%s%s?max=2", f.apiBase, f.User)
    
    // Log the URL being requested
    log.Printf("Fetching games for user: %s from URL: %s\n", f.User, url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Printf("Error creating request: %v\n", err)
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-ndjson")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request: %v\n", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v\n", err)
        return nil, err
    }

    // Log the number of games fetched (or if no games found)
    log.Printf("Received response with %d bytes of data\n", len(body))

    return parseGames(string(body)), nil
}

func parseGames(data string) []models.Game {
	var games []models.Game
	pgns := strings.Split(data, "\n\n\n") // games by double newlines

	log.Printf("Parsing %d games...\n", len(pgns))

	for _, pgn := range pgns {
		if strings.TrimSpace(pgn) == "" {
			continue
		}

		// Log the game being parsed
		log.Println("Parsing a new game...")

		game := chess.NewGame()
		if err := game.UnmarshalText([]byte(pgn)); err != nil {
			log.Printf("Error parsing PGN: %v\n", err)
			continue
		}

		headers := extractHeaders(pgn)

		// Log the parsed headers of the game
		log.Printf("Parsed game headers: %+v\n", headers)

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

	log.Printf("Parsed a total of %d games.\n", len(games))
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