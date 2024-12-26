package models

type Game struct {
    Event            string
    Site             string
    Date             string
    White            string
    Black            string
    Result           string
    UTCDate          string
    UTCTime          string
    WhiteElo         string
    BlackElo         string
    WhiteRatingDiff  string
    BlackRatingDiff  string
    Variant          string
    TimeControl      string
    ECO              string
    Termination      string
    Moves            string
}
