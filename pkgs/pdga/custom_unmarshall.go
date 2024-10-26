package pdga

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)


const dateTimeFormat = time.DateTime

type Time struct {
    time.Time
}

// UnmarshalJSON is a custom unmarshaler to get a time format from the API into a time.
func (ct *Time) UnmarshalJSON(b []byte) (err error) {
    s := strings.Trim(string(b), "\"")
    if s == "null" {
        ct.Time = time.Time{}
        return 
    }

    ct.Time, err = time.Parse(dateTimeFormat, s)

    return
}

type RoundDataUnmarshaler struct {
    RoundData []TournamentRoundData `json:"data"`
}


// UnmarshalJSON is a custom unmarshaler because the PDGA Live API will sometimes return different data.
// So we first check for a single item and unmarshal that, then we check for an array.
// We will always give a slice back for consistency.
func (d *RoundDataUnmarshaler) UnmarshalJSON(data []byte) error {
    var item TournamentRoundData

    singleItemErr := json.Unmarshal(data, &item)
    if singleItemErr == nil {
        d.RoundData = []TournamentRoundData{item}
        return nil
    }

    var items []TournamentRoundData
    sliceErr := json.Unmarshal(data, &items)
    if sliceErr != nil {
        return errors.Join(singleItemErr, sliceErr)
    }

    d.RoundData = items
    return nil
}

type SortRounds struct {
   Value string 
}

// SortRounds seems to be an interesting one.
// It most of the time comes back as an actual string most of the times, but sometimes just a 0.
// We first try to just unmarshal as a string, if that fails then just set as "0" for right now.
func (s *SortRounds) UnmarshalJSON(data []byte) error {
    var item string
    err := json.Unmarshal(data, &item) 
    if err != nil {
        s.Value = "0"
        return nil
    }

    s.Value = item
    return nil
}

// RoundScore holds the score value for a round for a player.
// this had to have a custom unmarshaller because it seems this comes back as a string on a few occassions.
type RoundScore int

func (score *RoundScore) UnmarshalJSON(data []byte) error {
    var result int
    err := json.Unmarshal(data, &result)
    if err == nil {
        score = (*RoundScore)(&result)
        return nil
    }

    var resultString string
    stringErr := json.Unmarshal(data, &resultString)
    if stringErr != nil {
        return err
    }
    
    result, err = strconv.Atoi(resultString)
    if err != nil {
        return err
    }
    
    score = (*RoundScore)(&result)
    return nil
}
