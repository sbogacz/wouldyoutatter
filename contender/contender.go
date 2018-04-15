package contender

// Contender is the model for the tattoo options
type Contender struct {
	Name        string
	Description string
	SVG         []byte
	Wins        int
	Losses      int
	Score       int
	isLoser     bool
}

// NewWinner creates a new "winning" contender
func NewWinner(name string) *Contender {
	return &Contender{
		Name:    name,
		isLoser: false,
	}
}

// NewLoser creates a new "losing" contender
func NewLoser(name string) *Contender {
	return &Contender{
		Name:    name,
		isLoser: true,
	}
}

// Matchup is the model for the head-to-head records
// between contenders
type Matchup struct {
	Contender1     string
	Contender2     string
	Contender1Wins int
	Contender2Wins int
}
