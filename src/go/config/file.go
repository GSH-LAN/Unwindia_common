package config

import (
	"context"
	"github.com/GSH-LAN/Unwindia_common/src/go/matchservice"
	"github.com/GSH-LAN/Unwindia_common/src/go/unwindiaError"
	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type ConfigFileImpl struct {
	ctx                context.Context
	currentConfig      *Config
	configFilename     string
	templatesDirectory string
	watcher            *fsnotify.Watcher
	lock               sync.RWMutex
}

func (c *ConfigFileImpl) GetGameServerTemplateForMatch(info matchservice.MatchInfo) (*GamerServerConfigTemplate, error) {
	gameName := info.Game
	gameValid := false
	var gsTemplate GamerServerConfigTemplate
	for game, cfg := range c.GetConfig().UnwindiaPteroConfig.Configs {
		if game == gameName {
			gameValid = true
			gsTemplate = cfg
			break
		}
	}

	if !gameValid {
		return nil, unwindiaError.NewInvalidGameError(gameName)
	}

	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
	}

	// TODO: make this shit reliable even with kinda broken configs
	// we now parse environments to replace custom variables and convert numeric values
	var newEnvironment = make(map[string]interface{})
	for envName, envValue := range gsTemplate.Environment {
		if val, ok := envValue.(string); ok {

			parsedEnvVar := strings.Builder{}
			err := template.Must(template.New("serverEnvironment").Funcs(funcs).Parse(val)).Execute(&parsedEnvVar, nil)
			if err != nil {
				log.Error().Err(err).Msg("Error parsing environment template")
			} else {
				val = parsedEnvVar.String()
			}

			if intValue, err := strconv.Atoi(val); err == nil {
				envValue = intValue
				log.Trace().Str("environment", envName).Int("value", intValue).Msg("Environemt is type int (parsed from string)")
			} else {
				envValue = val
				log.Trace().Str("environment", envName).Str("value", val).Msg("Environemt is type string")
			}
		}
		newEnvironment[envName] = envValue
	}

	gsTemplate.Environment = newEnvironment

	return &gsTemplate, nil
}

func (c *ConfigFileImpl) GetGameServerTemplate(gameName string) (*GamerServerConfigTemplate, error) {
	return c.GetGameServerTemplateForMatch(matchservice.MatchInfo{Game: gameName})
}

func NewConfigFile(ctx context.Context, filename, templatesDirectory string) (ConfigClient, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// fileWatcher for config-file
	err = watcher.Add(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("Error adding config file watcher")
	}

	// fileWatcher for template-files
	err = watcher.Add(templatesDirectory)
	if err != nil {
		log.Fatal().Err(err).Msg("Error adding template file watcher")
	}

	go func() {
		for {
			// wait for context closed
			select {
			case <-ctx.Done():
				_ = watcher.Close()
			}
		}
	}()

	cfg := ConfigFileImpl{
		ctx:                ctx,
		configFilename:     filename,
		templatesDirectory: templatesDirectory,
		watcher:            watcher,
	}

	_, err = cfg.loadConfig()
	if err != nil {
		return nil, err
	}

	go cfg.startFileWatcher()

	return &cfg, nil
}
func (c *ConfigFileImpl) GetConfig() *Config {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.currentConfig == nil {
		_, err := c.loadConfig()
		if err != nil {
			log.Error().Err(err).Msg("Error loading config")
		}
	}

	return c.currentConfig
}

func (c *ConfigFileImpl) loadConfig() (*Config, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var cfg Config

	configFile, err := os.ReadFile(c.configFilename)
	if err != nil {
		log.Error().Err(err).Str("filename", c.configFilename).Msg("Error loading config file")
		return nil, err
	}

	if err = jsoniter.Unmarshal(configFile, &cfg); err != nil {
		return c.currentConfig, err
	}

	log.Info().Msgf("Loaded config from file: %+v", cfg)

	templates := make(map[string]string)

	if files, err := os.ReadDir(c.templatesDirectory); err != nil {
		log.Warn().Err(err).Msg("Error loading templates")
	} else {
		for _, file := range files {
			if filecontent, err := os.ReadFile(path.Join(c.templatesDirectory, file.Name())); err != nil {
				log.Warn().Err(err).Str("filename", file.Name()).Msg("Error loading content of file")
			} else {
				templates[file.Name()] = string(filecontent)
			}
		}
		log.Info().Msgf("Loaded templates: %+v", templates)
		cfg.Templates = templates
	}

	if !reflect.DeepEqual(cfg, Config{}) || c.currentConfig == nil {
		c.currentConfig = &cfg
	}

	return c.currentConfig, nil
}

func (c *ConfigFileImpl) startFileWatcher() {
	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				return
			}
			log.Debug().Msgf("Config file watcher event: %+v", event)
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) || event.Has(fsnotify.Chmod) {
				log.Info().Str("event", event.Name).Msgf("modified config file")
				_, err := c.loadConfig()
				if err != nil {
					log.Error().Err(err).Msg("Error loading config from file")
				}
			}
		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("Filewatcher encountered error")
		}
	}
}
