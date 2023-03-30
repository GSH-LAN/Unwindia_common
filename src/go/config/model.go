package config

import (
	"errors"
	"github.com/GSH-LAN/Unwindia_common/src/go/matchservice"
	"github.com/GSH-LAN/Unwindia_common/src/go/messagebroker"
	jsoniter "github.com/json-iterator/go"
	"github.com/parkervcp/crocgodyl"
	"time"
)

type ConfigClient interface {
	GetConfig() *Config
	GetGameServerTemplate(gameName string) (*GamerServerConfigTemplate, error)
	GetGameServerTemplateForMatch(info matchservice.MatchInfo) (*GamerServerConfigTemplate, error)
}

type Config struct {
	Templates            map[string]string          `json:"templates,omitempty"`
	CmsConfig            CmsConfig                  `json:"cmsConfig"`
	UpdateDotlanOnEvents []messagebroker.MatchEvent `json:"updateDotlanOnEvents"`
	UnwindiaPteroConfig  `json:"pterodactyl"`
}

type CmsConfig struct {
	UserId           uint   `json:"dotlanUserId"`
	TournamentFilter string `json:"dotlanTournamentFilter"`
	DefaultGame      string `json:"defaultGame"`
}

type UnwindiaPteroConfig struct {
	Configs map[string]GamerServerConfigTemplate `json:"configs"` // configs is a map which contains the game name as key with belonging template
}

type GamerServerConfigTemplate struct {
	UserId                  int                    `json:"userId"`
	LocationId              int                    `json:"locationId"`
	NestId                  int                    `json:"nestId"`
	ServerNamePrefix        string                 `json:"serverNamePrefix"`
	ServerNameGOTVPrefix    string                 `json:"serverNameGOTVPrefix"`
	DefaultStartup          string                 `json:"defaultStartup"`
	DefaultDockerImage      string                 `json:"defaultDockerImage"`
	DefaultServerPassword   string                 `json:"defaultServerPassword"`
	DefaultRconPassword     string                 `json:"defaultRconPassword"`
	EggId                   int                    `json:"eggId"`
	Limits                  crocgodyl.Limits       `json:"limits"`      // Limits which will be set for new servers.
	ForceLimits             bool                   `json:"forceLimits"` // if true, server which does not meet Limits settings will be deleted. If false, existing suspended servers will be reused, no matter of matching Limits
	Environment             map[string]interface{} `json:"environment"` // Environment settings for a game. Can be values or matching properties of a gameserver match thingy object in go-template like format // TODO: set correct object description
	TvSlots                 int                    `json:"tvSlots"`
	DeleteAfterDuration     Duration               `json:"deleteAfterDuration"`
	ServerReadyRconCommand  string                 `json:"serverReadyRconCommand"`
	ServerReadyRconWaitTime Duration               `json:"serverReadyRconWaitTime"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := jsoniter.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
