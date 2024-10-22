package pdga

import (
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
}

func NewClient(options ...Option) *Client {
    client := &Client{
        httpClient: http.DefaultClient,
    }
    
    for _, option := range options {
        option(client)
    }

    return client
}

func (c *Client) FetchTournamentRound(tournamentId, roundNumber, int, division Division) {

}
