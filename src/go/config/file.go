package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/GSH-LAN/Unwindia_common/src/go/logger"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

type ConfigFileImpl struct {
	ctx                context.Context
	currentConfig      *Config
	configFilename     string
	templatesDirectory string
	watcher            *fsnotify.Watcher
}

func NewConfigFile(ctx context.Context, filename, templatesDirectory string) (ConfigClient, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			// wait for context closed
			if ctx.Err() != nil {
				watcher.Close()
			}
		}
	}()

	cfg := &ConfigFileImpl{
		ctx:                ctx,
		configFilename:     filename,
		templatesDirectory: templatesDirectory,
		watcher:            watcher,
	}

	_, err = cfg.loadConfig()
	if err != nil {
		return nil, err
	}

	cfg.startFileWatcher()

	return cfg, nil
}
func (c *ConfigFileImpl) GetConfig() *Config {

	if c.currentConfig == nil {
		_, err := c.loadConfig()
		if err != nil {
			log.Error("Error loading config: %+v", err)
		}
	}

	return c.currentConfig
}

func (c *ConfigFileImpl) loadConfig() (*Config, error) {
	var cfg Config

	configFile, err := ioutil.ReadFile(c.configFilename)
	if err != nil {
		log.Errorf("Error loading config file %s: %+v", c.configFilename, err)
		return nil, err
	}

	json.Unmarshal(configFile, &cfg)

	log.Infof("Loaded config from file: %+v", cfg)

	templates := make(map[string]string)

	if files, err := ioutil.ReadDir(c.templatesDirectory); err != nil {
		log.Warnf("Error loading templates: %+v", err)
	} else {
		for _, file := range files {
			if filecontent, err := ioutil.ReadFile(path.Join(c.templatesDirectory, file.Name())); err != nil {
				log.Warnf("Error loading content of file %s: %+v", file.Name(), err)
			} else {
				templates[file.Name()] = string(filecontent)
			}
		}
		log.Infof("Loaded templates: %+v", templates)
		cfg.Templates = templates
	}

	c.currentConfig = &cfg

	return c.currentConfig, nil
}

func (c *ConfigFileImpl) startFileWatcher() {
	go func() {
		for {
			select {
			case event, ok := <-c.watcher.Events:
				if !ok {
					return
				}
				log.Infof("Config file watcher event: %+v", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Infof("modified config file: %+v", event.Name)
					_, err := c.loadConfig()
					if err != nil {
						log.Errorf("Error loading configfrom file : %+v", err)
					}
				}
			case err, ok := <-c.watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("error: %+v", err)
			}
		}
	}()
}
