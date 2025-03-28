package madden

// ExportData represents the structure of data exported from the Madden Companion App
// This is a basic structure and might need to be expanded based on actual data format
type ExportData struct {
	ExportType string                 `json:"exportType"`
	Timestamp  string                 `json:"timestamp,omitempty"`
	LeagueID   string                 `json:"leagueId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// Player represents a player in Madden
type Player struct {
	PlayerID      int    `json:"playerId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	JerseyNum     int    `json:"jerseyNum"`
	Position      string `json:"position"`
	TeamID        int    `json:"teamId"`
	Age           int    `json:"age"`
	Height        int    `json:"height"`
	Weight        int    `json:"weight"`
	YearsPro      int    `json:"yearsPro"`
	PlayerBestOvr int    `json:"playerBestOvr"`
	// Additional attributes can be added as needed
}

// Team represents a team in Madden
type Team struct {
	TeamID      int    `json:"teamId"`
	DisplayName string `json:"displayName"`
	TeamOvr     int    `json:"teamOvr"`
	City        string `json:"city"`
	Nickname    string `json:"nickname"`
	DefScheme   string `json:"defScheme"`
	OffScheme   string `json:"offScheme"`
	// Additional attributes can be added as needed
}

// LeagueInfo represents information about a Madden league
type LeagueInfo struct {
	LeagueID   string `json:"leagueId"`
	LeagueName string `json:"leagueName"`
	SeasonYear int    `json:"seasonYear"`
	SeasonWeek int    `json:"seasonWeek"`
	StageIndex int    `json:"stageIndex"`
	StageWeek  int    `json:"stageWeek"`
}
