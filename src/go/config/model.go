package config

import "github.com/GSH-LAN/Unwindia_common/src/go/messagebroker"

type ConfigClient interface {
	GetConfig() *Config
}

type Config struct {
	Templates            map[string]string      `json:"templates,omitempty"`
	CmsConfig            CmsConfig              `json:"cmsConfig"`
	UpdateDotlanOnEvents []messagebroker.Events `json:"updateDotlanOnEvents"`
}

type CmsConfig struct {
	UserId           uint   `json:"dotlanUserId"`
	TournamentFilter string `json:"dotlanTournamentFilter"`
	DefaultGame      string `json:"defaultGame"`
}
