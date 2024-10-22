package pdga

import "encoding/json"

func UnmarshalTournamentRoundData(data []byte) (TournamentRoundData, error) {
	var r TournamentRoundData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TournamentRoundData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TournamentRoundData struct {
	Data Data   `json:"data"`
	Hash string `json:"hash"`
}

type Data struct {
	Pool        string   `json:"pool"`
	Layouts     []Layout `json:"layouts"`
	Division    Division `json:"division"`
	LiveRoundID int64    `json:"live_round_id"`
	ID          int64    `json:"id"`
	ShotgunTime string   `json:"shotgun_time"`
	TeeTimes    bool     `json:"tee_times"`
	Holes       []Hole   `json:"holes"`
	Scores      []Score  `json:"scores"`
}

type Hole struct {
	Hole        string      `json:"Hole"`
	HoleOrdinal int64       `json:"HoleOrdinal"`
	Label       string      `json:"Label"`
	Par         int64       `json:"Par"`
	Length      int64       `json:"Length"`
	Units       interface{} `json:"Units"`
	Accuracy    interface{} `json:"Accuracy"`
	Ordinal     int64       `json:"Ordinal"`
}

type Layout struct {
	LayoutID        int64       `json:"LayoutID"`
	CourseID        int64       `json:"CourseID"`
	CourseName      string      `json:"CourseName"`
	TournID         int64       `json:"TournID"`
	Name            string      `json:"Name"`
	Holes           int64       `json:"Holes"`
	Par             int64       `json:"Par"`
	Length          int64       `json:"Length"`
	Units           string      `json:"Units"`
	Accuracy        string      `json:"Accuracy"`
	Notes           string      `json:"Notes"`
	H1              int64       `json:"H1"`
	H2              int64       `json:"H2"`
	H3              int64       `json:"H3"`
	H4              int64       `json:"H4"`
	H5              int64       `json:"H5"`
	H6              int64       `json:"H6"`
	H7              int64       `json:"H7"`
	H8              int64       `json:"H8"`
	H9              int64       `json:"H9"`
	H10             int64       `json:"H10"`
	H11             int64       `json:"H11"`
	H12             int64       `json:"H12"`
	H13             int64       `json:"H13"`
	H14             int64       `json:"H14"`
	H15             int64       `json:"H15"`
	H16             int64       `json:"H16"`
	H17             int64       `json:"H17"`
	H18             int64       `json:"H18"`
	H19             int64       `json:"H19"`
	H20             int64       `json:"H20"`
	H21             int64       `json:"H21"`
	H22             int64       `json:"H22"`
	H23             int64       `json:"H23"`
	H24             int64       `json:"H24"`
	H25             int64       `json:"H25"`
	H26             int64       `json:"H26"`
	H27             int64       `json:"H27"`
	H28             int64       `json:"H28"`
	H29             int64       `json:"H29"`
	H30             int64       `json:"H30"`
	H31             int64       `json:"H31"`
	H32             int64       `json:"H32"`
	H33             int64       `json:"H33"`
	H34             int64       `json:"H34"`
	H35             int64       `json:"H35"`
	H36             int64       `json:"H36"`
	UpdateDate      Time   `json:"UpdateDate"`
	Detail          []Hole      `json:"Detail"`
}

type Score struct {
	ResultID            int64             `json:"ResultID"`
	RoundID             int64             `json:"RoundID"`
	ScoreID             int64             `json:"ScoreID"`
	FirstName           string            `json:"FirstName"`
	LastName            string            `json:"LastName"`
	Name                string            `json:"Name"`
	AvatarURL           string            `json:"AvatarURL"`
	City                string            `json:"City"`
	Country             string           `json:"Country"`
	Nationality         string           `json:"Nationality"`
	StateProv           string           `json:"StateProv"`
	PDGANum             int64             `json:"PDGANum"`
	HasPDGANum          int64             `json:"HasPDGANum"`
	Rating              int64             `json:"Rating"`
	Division            Division          `json:"Division"`
	Pool                string            `json:"Pool"`
	Round               int64             `json:"Round"`
	Authoritative       int64             `json:"Authoritative"`
	ScorecardUpdatedAt  Time         `json:"ScorecardUpdatedAt"`
	WonPlayoff          WonPlayoff        `json:"WonPlayoff"`
	Prize               string            `json:"Prize"`
	PrevRounds          int64             `json:"PrevRounds"`
	RoundStatus         string       `json:"RoundStatus"`
	Holes               int64             `json:"Holes"`
	Par                 int64             `json:"Par"`
	LayoutID            int64             `json:"LayoutID"`
	GrandTotal          int64             `json:"GrandTotal"`
	CardNum             int64             `json:"CardNum"`
	TeeTime             string            `json:"TeeTime"`
	TeeStart            string            `json:"TeeStart"`
	HasGroupAssignment  int64             `json:"HasGroupAssignment"`
	PlayedPreviousRound int64             `json:"PlayedPreviousRound"`
	HasRoundScore       int64             `json:"HasRoundScore"`
	UpdateDate          Time         `json:"UpdateDate"`
	Played              int64             `json:"Played"`
	Completed           int64             `json:"Completed"`
	RoundStarted        int64             `json:"RoundStarted"`
	PrevRndTotal        int64             `json:"PrevRndTotal"`
	RoundScore          int64             `json:"RoundScore"`
	SubTotal            int64             `json:"SubTotal"`
	RoundtoPar          int64             `json:"RoundtoPar"`
	ToPar               int64             `json:"ToPar"`
	Scores              string            `json:"Scores"`
	SortScores          string            `json:"SortScores"`
	Pars                string            `json:"Pars"`
	Rounds              string            `json:"Rounds"`
	SortRounds          string            `json:"SortRounds"`
	RoundRating         int64             `json:"RoundRating"`
	PreviousPlace       int64             `json:"PreviousPlace"`
	FullLocation        string            `json:"FullLocation"`
	ShortName           string            `json:"ShortName"`
	ProfileURL          string            `json:"ProfileURL"`
	ParThruRound        int64             `json:"ParThruRound"`
	RoundPool           string            `json:"RoundPool"`
	TeeTimeSort         string            `json:"TeeTimeSort"`
	RunningPlace        int64             `json:"RunningPlace"`
	Tied                bool              `json:"Tied"`
	HoleScores          []string          `json:"HoleScores"`
}

type WonPlayoff string

const (
	No WonPlayoff = "no"
)
