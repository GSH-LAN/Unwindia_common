package config

type ConfigClientImpl struct {
	currentConfig *Config
}

func NewConfigClient() (ConfigClient, error) {
	return ConfigClientImpl{}, nil
}

// TODO: Implement client for config service
func (c ConfigClientImpl) GetConfig() *Config {
	if c.currentConfig == nil {
		cfg := Config{
			Templates: map[string]string{
				"CMS_FORUM_POST": "Hallo, ich bin der Turnierbot.\n\nIch werde euch durch diese Begegnung begleiten. Sobald beide Teams bereit sind, werde ich einen Gameserver für euch generieren und die Zugangsdaten in diesem Kommentar hinterlegen.\n\n",
			},
			CmsConfig: CmsConfig{
				UserId:           3,
				TournamentFilter: "teventid = 1 and tgameserver = 1",
				DefaultGame:      "csgo",
			},
		}

		c.currentConfig = &cfg
	}

	panic("not implemented")
}
