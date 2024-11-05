package pdga

import "encoding/json"

// TODO: I need to turn into a struct that holds round information
type RoundValue struct {
    RoundKey string
    Data RoundsData
}

func (r *RoundValue) UnmarshallJSON(data []byte) error {
    roundMap := make(map[string]RoundsData)
    if err := json.Unmarshal(data, &roundMap); err != nil {
       return err 
    }
    
    return nil
}
