package pdga

import (
	"encoding/json"
	"strconv"
)

func UnmarshalTournamentInfo(data []byte) (TournamentInfo, error) {
	var r TournamentInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TournamentInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TournamentInfo struct {
	Data TournamentData `json:"data"`
	Hash string         `json:"hash"`
}

type TournamentData struct {
	DateRange             string                `json:"DateRange"`
	EndDate               string                `json:"EndDate"`
	Country               string                `json:"Country"`
	Location              string                `json:"Location"`
	LocationShort         string                `json:"LocationShort"`
	Divisions             []TournamentDivision  `json:"Divisions"`
	Finals                string                `json:"Finals"`
	HighestCompletedRound int64                 `json:"HighestCompletedRound"`
	LatestRound           string                `json:"LatestRound"`
	FinalRound            int64                 `json:"FinalRound"`
	Layouts               []LayoutInfo          `json:"Layouts"`
	LayoutIDs             []LayoutID            `json:"LayoutIDs"`
	Name                  string                `json:"Name"`
	NameHTML              string                `json:"NameHtml"`
	MultiLineName         MultiLineName         `json:"MultiLineName"`
	MultiLineNameHTML     string                `json:"MultiLineNameHtml"`
	RawTier               string                `json:"RawTier"`
	Rounds                int64                 `json:"Rounds"`
	RoundsList            map[string]RoundsList `json:"RoundsList"`
	Semis                 string                `json:"Semis"`
	SimpleName            string                `json:"SimpleName"`
	StartDate             string                `json:"StartDate"`
	TDName                string                `json:"TDName"`
	TDPDGANum             int64                 `json:"TDPDGANum"`
	Tier                  string                `json:"Tier"`
	FormattedTier         string                `json:"FormattedTier"`
	FormattedLongTier     string                `json:"FormattedLongTier"`
	TierPro               string                `json:"TierPro"`
	TierAm                string                `json:"TierAm"`
	TierX                 bool                  `json:"TierX"`
	TotalPlayers          int64                 `json:"TotalPlayers"`
	TournamentID          string                `json:"TournamentId"`
	International         bool                  `json:"International"`
	FullStats             bool                  `json:"FullStats"`
	AdditionalEventInfo   AdditionalEventInfo   `json:"AdditionalEventInfo"`
}

type AdditionalEventInfo struct {
	RoundsList   map[string]RoundsList `json:"RoundsList"`
	BroadcastURL string                `json:"BroadcastUrl"`
}

type RoundsList struct {
	Number           int64  `json:"Number"`
	Label            string `json:"Label"`
	LabelAbbreviated string `json:"LabelAbbreviated"`
	Date             string `json:"Date"`
	DateAbbreviated  string `json:"DateAbbreviated"`
	ShowDate         bool   `json:"ShowDate"`
}

type TournamentDivision struct {
	DivisionID        int64            `json:"DivisionID"`
	Division          string           `json:"Division"`
	DivisionName      string           `json:"DivisionName"`
	Players           int64            `json:"Players"`
	LayoutAssignments map[string]int64 `json:"LayoutAssignments"`
	IsPro             bool             `json:"IsPro"`
	ShortName         string           `json:"ShortName"`
	AbbreviatedName   string           `json:"AbbreviatedName"`
	LatestRound       string           `json:"LatestRound"`
}

type LayoutID struct {
	The1     int64  `json:"1"`
	The2     int64  `json:"2"`
	The3     int64  `json:"3"`
	The4     int64  `json:"4"`
	The5     int64  `json:"5"`
	The6     int64  `json:"6"`
	The7     int64  `json:"7"`
	The8     int64  `json:"8"`
	The9     int64  `json:"9"`
	The10    int64  `json:"10"`
	The11    int64  `json:"11"`
	The12    int64  `json:"12"`
	The13    int64  `json:"13"`
	Division string `json:"Division"`
	Pool     string `json:"Pool"`
}

type LayoutInfo struct {
	LayoutID   int64    `json:"LayoutID"`
	CourseID   int64    `json:"CourseID"`
	CourseName string   `json:"CourseName"`
	Name       string   `json:"Name"`
	Holes      int64    `json:"Holes"`
	Par        int64    `json:"Par"`
	Length     int64    `json:"Length"`
	Units      string   `json:"Units"`
	Accuracy   string   `json:"Accuracy"`
	Details    []Detail `json:"Details"`
}

type Detail struct {
	Hole   string `json:"Hole"`
	Label  string `json:"Label"`
	Par    int64  `json:"Par"`
	Length *int64 `json:"Length"`
}

type MultiLineName struct {
	Pre  string `json:"pre"`
	Main string `json:"main"`
	Post string `json:"post"`
}

func (td TournamentData) NumberRounds() int {
	numberRounds := len(td.RoundsList)
	playoffRound := 0
	for _, value := range td.RoundsList {
		if value.Label == "Playoff" {
			playoffRound = int(value.Number)
		}
	}

	if td.HighestCompletedRound != int64(playoffRound) {
		numberRounds--
	}

	return numberRounds
}

// Attempts to get the tournament id as an int.
// TODO: See about tourning this into a separate type to compare later?
func (t TournamentInfo) IdAsInt() (int, error) {
	return strconv.Atoi(t.Data.TournamentID)
}
