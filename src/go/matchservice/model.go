package matchservice

// Player represents a player which is a part of a team
type Player struct {
	Id             string // The id of the player in the external tournament system
	Name           string // Name of the player in the external tournament system
	GameProviderId string // Id of an game provider, e.g. SteamId for CS:GO
	Picture        []byte `json:"picture,omitempty"` // Picture of the player, optional
}

// Team represents a team which is a part of a match
type Team struct {
	Id      string   // The id of the team in the external tournament system
	Name    string   // Name of the team in the external tournament system
	Players []Player // List of players in the team
	Picture []byte   `json:"picture,omitempty"` // Picture of the team, optional
}

type MatchInfo struct {
	Id                 string // Unwindia Match ID
	MsId               string // MatchService Match-ID
	Team1              Team   // Team 1
	Team2              Team   // Team 2
	PlayerAmount       uint   // Amount of players in the match
	Game               string // Game name
	Map                string // Map name
	ServerAddress      string // Server address
	ServerPassword     string // Server password
	ServerPasswordMgmt string // Server password for management
	ServerTvAddress    string // Server TV address
	ServerTvPassword   string // Server TV password
}
