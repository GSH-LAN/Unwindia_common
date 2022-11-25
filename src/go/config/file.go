package config

import (
	"context"
	"github.com/fsnotify/fsnotify"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"reflect"
	"sync"
)

type ConfigFileImpl struct {
	ctx                context.Context
	currentConfig      *Config
	configFilename     string
	templatesDirectory string
	watcher            *fsnotify.Watcher
	lock               sync.RWMutex
}

func NewConfigFile(ctx context.Context, filename, templatesDirectory string) (ConfigClient, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("Error adding config file watcher")
	}

	go func() {
		for {
			// wait for context closed
			if ctx.Err() != nil {
				watcher.Close()
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
			if event.Op&fsnotify.Write == fsnotify.Write {
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
