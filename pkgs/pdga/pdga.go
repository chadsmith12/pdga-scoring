package pdga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
    tournamentEndpoint = "live_results_fetch_event"
)

var defaultClient = NewClient()

func DownloadTournament(basePath string, tournamentId int) error {
    tournament, err := defaultClient.FetchTournamentInfo(context.Background(), tournamentId)
    if err != nil {
        return err
    }
    
    data, err := json.Marshal(tournament)
    if err != nil {
        return err
    }
    
    tournamentPath := path.Join(basePath, fmt.Sprintf("%d", tournamentId), "tournament.json")
    if err := os.MkdirAll(filepath.Dir(tournamentPath), 0755); err != nil {
        return err
    }

    return os.WriteFile(tournamentPath, data, 0755)
}

func DownloadRoundData(basePath string, tournamentId, roundNumber int) error {
    mpoRound, err := defaultClient.FetchTournamentRound(context.Background(), tournamentId, roundNumber, Mpo)
    if err != nil {
        return err
    }
    mpoData, err := json.Marshal(mpoRound)
    if err != nil {
        return err
    }

    if err := writeDivisionRound(basePath, tournamentId, roundNumber, Mpo, mpoData); err != nil {
        return err
    }


    fpoRound, err := defaultClient.FetchTournamentRound(context.Background(), tournamentId, roundNumber, Fpo)
    if err != nil {
        return err
    }

    fpoData, err := json.Marshal(fpoRound)
    if err != nil {
        return err
    }

    if err := writeDivisionRound(basePath, tournamentId, roundNumber, Fpo, fpoData); err != nil {
        return err
    }

    return nil
}

func writeDivisionRound(basePath string, tournamentId, roundNumber int, division Division, data []byte) error {
    path := path.Join(basePath, fmt.Sprintf("%d", tournamentId), fmt.Sprintf("%s-%d.json", division, roundNumber))
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return err
    }
    
    return os.WriteFile(path, data, 0755)
}

type Client struct {
    httpClient *http.Client
    baseUrl string
}

func NewClient(options ...Option) *Client {
    client := &Client{
        httpClient: http.DefaultClient,
        baseUrl: baseUrl,
    }
    
    for _, option := range options {
        option(client)
    }

    return client
}

func WithClient(client *http.Client) Option {
    return func(c *Client) {
        c.httpClient = client
    }
}

func WithBaseUrl(baseUrl string) Option {
    return func(c *Client) {
       c.baseUrl = baseUrl 
    }
}

func (c *Client) FetchTournamentInfo(ctx context.Context, tournamentId int) (TournamentInfo, error) {
    tournamentUrl := c.createTournamentUrl(tournamentId)
    req, err := http.NewRequestWithContext(ctx, "GET", tournamentUrl, nil)
    if err != nil {
        return TournamentInfo{}, err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return TournamentInfo{}, err
    }
    if resp.StatusCode != 200 {
        return TournamentInfo{}, fmt.Errorf("tournament info returned status code of %d", resp.StatusCode)
    }

    body := resp.Body
    defer body.Close()

    var tournamentData TournamentInfo
    err = json.NewDecoder(body).Decode(&tournamentData)
   
    return tournamentData, err
}

func (c *Client) FetchTournamentRound(context context.Context, tournamentId, roundNumber int, division Division) (TournamentRoundResponse, error) {
    roundUrl := c.createTournamentRoundUrl(tournamentId, roundNumber, division)
    req, err := http.NewRequestWithContext(context, "GET", roundUrl, nil)
    if err != nil {
        return TournamentRoundResponse{}, err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return TournamentRoundResponse{}, err
    }

    if resp.StatusCode != 200 {
        return TournamentRoundResponse{}, errors.New(fmt.Sprintf("round results returned status code of %d", resp.StatusCode))
    }

    body := resp.Body
    defer body.Close()
    
    var tournamentRoundData TournamentRoundResponse
    err = json.NewDecoder(body).Decode(&tournamentRoundData)
    if err != nil {
        return TournamentRoundResponse{}, err
    }

    return tournamentRoundData, nil
}

func (c *Client) createTournamentRoundUrl(tournamentId, roundNumber int, division Division) string {
    return fmt.Sprintf("%s/%s?TournID=%d&Division=%s&Round=%d", c.baseUrl, roundEndpoint, tournamentId, division, roundNumber)
}

func (c *Client) createTournamentUrl(tournamentId int) string {
    return fmt.Sprintf("%s/%s?TournID=%d", c.baseUrl, tournamentEndpoint, tournamentId)
}
