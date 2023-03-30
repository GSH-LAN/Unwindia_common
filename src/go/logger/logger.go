// Package logger configures the zerolog logger
package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {
	log.Logger = log.Hook(LineInfoHook{})

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).Level(zerolog.InfoLevel).With().Timestamp().Logger()

}

func SetLogLevel(level string) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Error().Err(err).Str("level", level).Msg("Error parsing log level")
		return err
	}

	log.Logger = log.Level(lvl)
	log.Info().Str("loglevel", level).Msgf("Changed loglevel")

	return nil
}

type LineInfoHook struct{}

func (h LineInfoHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	_, file, line, ok := runtime.Caller(0)
	if ok {
		e.Str("line", fmt.Sprintf("%s:%d", file, line))
	}
}
