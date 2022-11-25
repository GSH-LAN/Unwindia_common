package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func init() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
}

func Set(level string) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Error().Err(err).Str("level", level).Msg("Error parsing log level")
		return err
	}

	log.Logger = log.Level(lvl)
	log.Logger.Info().Str("level", level).Msg("Changed log level")

	return nil
}
