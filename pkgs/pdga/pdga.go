package pdga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Option func(*Client)
type Division string

const (
    Mpo Division = "MPO"
    Fpo Division = "FPO"
)

const ( 
    baseUrl = "https://www.pdga.com/apps/tournament/live-api"
    roundEndpoint = "live_results_fetch_round"
)


const dateTimeFormat = time.DateTime

type Time struct {
    time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) (err error) {
    s := strings.Trim(string(b), "\"")
    if s == "null" {
        ct.Time = time.Time{}
        return 
    }

    ct.Time, err = time.Parse(dateTimeFormat, s)

    return
}

type Client struct {
    httpClient *http.Client
    context context.Context
}

func NewClient(options ...Option) *Client {
    client := &Client{
        httpClient: http.DefaultClient,
        context: context.Background(),
    }
    
    for _, option := range options {
        option(client)
    }

    return client
}

func (c *Client) FetchTournamentRound(tournamentId, roundNumber int, division Division) (TournamentRoundData, error) {
    roundUrl := createTournamentRoundUrl(tournamentId, roundNumber, division)
    req, err := http.NewRequestWithContext(c.context, "GET", roundUrl, nil)
    if err != nil {
        return TournamentRoundData{}, err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return TournamentRoundData{}, err
    }

    if resp.StatusCode != 200 {
        return TournamentRoundData{}, errors.New(fmt.Sprintf("round results returned status code of %d", resp.StatusCode))
    }

    body := resp.Body
    defer body.Close()
    
    var tournamentRoundData TournamentRoundData
    err = json.NewDecoder(body).Decode(&tournamentRoundData)
    if err != nil {
        return TournamentRoundData{}, err
    }

    return tournamentRoundData, nil
}

func createTournamentRoundUrl(tournamentId, roundNumber int, division Division) string {
    return fmt.Sprintf("%s/%s?TournID=%d&Division=%s&Round=%d", baseUrl, roundEndpoint, tournamentId, division, roundNumber)
}
