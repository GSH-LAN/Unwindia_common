package config

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/GSH-LAN/Unwindia_common/src/go/logger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

type ConfigFileImpl struct {
	currentConfig      *Config
	configFilename     string
	templatesDirectory string
}

func NewConfigFile(filename, templatesDirectory string) (ConfigClient, error) {
	cfg := &ConfigFileImpl{
		configFilename:     filename,
		templatesDirectory: templatesDirectory,
	}

	cfg.loadConfig()

	return cfg, nil
}
func (c *ConfigFileImpl) GetConfig() *Config {

	if c.currentConfig == nil {
		_, err := c.loadConfig()
		if err != nil {
			log.Error("Eoor loading config: %+v", err)
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
