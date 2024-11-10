package fantasy

import (
	"encoding/json"
	"io"
)

type Teams []Team

func UnmarshalTeams(data []byte) (Teams, error) {
	var r Teams
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Teams) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func LoadTeams(reader io.Reader) (Teams, error) {
    data, err := io.ReadAll(reader) 
    if err != nil {
        return []Team{}, err
    }

    return UnmarshalTeams(data)
}

type Team struct {
	Name string   `json:"Name"`
	Team TeamPlayers `json:"Team"`
}

type TeamPlayers struct {
	Players []int64 `json:"Players"`
	Bench   []int64 `json:"Bench"`
}
