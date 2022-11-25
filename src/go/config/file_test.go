package config

import (
	"context"
	"github.com/GSH-LAN/Unwindia_common/src/go/messagebroker"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

const configTestJson = `{
  "cmsConfig":{
    "dotlanUserId": 1,
    "dotlanTournamentFilter": "teventid = 1",
    "defaultGame": "csgo"
  },
  "updateDotlanOnEvents": ["UNWINDIA_MATCH_NEW"]
}`

var testConfig = Config{
	Templates: nil,
	CmsConfig: CmsConfig{
		UserId:           1,
		TournamentFilter: "teventid = 1",
		DefaultGame:      "csgo",
	},
	UpdateDotlanOnEvents: []messagebroker.MatchEvent{messagebroker.UNWINDIA_MATCH_NEW},
}

const configTestJsonUpdated = `{
  "cmsConfig":{
    "dotlanUserId": 2,
    "dotlanTournamentFilter": "teventid = 2",
    "defaultGame": "csgo"
  },
  "updateDotlanOnEvents": ["UNWINDIA_MATCH_NEW"]
}`

var testConfigUpdated = Config{
	Templates: nil,
	CmsConfig: CmsConfig{
		UserId:           2,
		TournamentFilter: "teventid = 2",
		DefaultGame:      "csgo",
	},
	UpdateDotlanOnEvents: []messagebroker.MatchEvent{messagebroker.UNWINDIA_MATCH_NEW},
}

const configTestJsonBroken = `{
  "cmsConfig":{
    "dotlanUserId": 1,
    "dotlanTournamentFilter": "teventid = 1",
    "defaultGame": "csgo"
  }
  "updateDotlanOnEvents": ["UNWINDIA_CMS_CONTEST_NEW"]
}`

var configTestFile *os.File

func TestMain(m *testing.M) {
	// Setup temp file for config
	file, err := ioutil.TempFile("", "config.*.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(file.Name())

	configTestFile = file

	_, err2 := configTestFile.WriteString(configTestJson)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Run the other tests
	os.Exit(m.Run())
}

func TestConfigFileImpl_GetConfig(t *testing.T) {
	a := assert.New(t)
	type fields struct {
		configFilename     string
		templatesDirectory string
	}
	tests := []struct {
		name   string
		fields fields
		want   Config
	}{
		{
			name: "load config",
			fields: fields{
				configFilename:     configTestFile.Name(),
				templatesDirectory: "",
			},
			want: testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewConfigFile(context.Background(), tt.fields.configFilename, tt.fields.templatesDirectory)
			a.NoError(err)
			if got := c.GetConfig(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileImpl_ReloadConfigOnChange(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	type fields struct {
		changedConfigJson string
	}

	tests := []struct {
		name   string
		fields fields
		want   Config
	}{
		{
			name: "no change",
			fields: fields{
				changedConfigJson: configTestJson,
			},
			want: testConfig,
		},
		{
			name: "updated config",
			fields: fields{
				changedConfigJson: configTestJsonUpdated,
			},
			want: testConfigUpdated,
		},
		{
			name: "broken json",
			fields: fields{
				changedConfigJson: configTestJsonBroken,
			},
			want: testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, err := ioutil.TempFile("", "config.*.json")
			if err != nil {
				log.Fatal(err)
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					log.Fatal(err)
				}
			}(file.Name())

			_, err = file.WriteString(configTestJson)
			a.NoError(err)

			c, err := NewConfigFile(context.Background(), file.Name(), "")
			a.NoError(err)

			c.GetConfig()

			err = file.Truncate(0)
			a.NoError(err)

			_, err = file.Seek(0, 0)
			a.NoError(err)

			_, err = file.WriteString(tt.fields.changedConfigJson)
			a.NoError(err)

			time.Sleep(time.Millisecond * 100)

			if got := c.GetConfig(); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
