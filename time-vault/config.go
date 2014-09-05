package timevault

import (
	"errors"

	"appengine"
	"appengine/datastore"
)

type Config struct {
	Active            bool   `datastore:"active"             json:"active"`
	TwilioURL         string `datastore:"twilioURL"          json:"twilioURL"`
	TwilioMessagePath string `datastore:"twilioMessagePath"  json:"twilioMessagePath"`
	TwilioSID         string `datastore:"twilioSID"          json:"twilioSID"`
	TwilioToken       string `datastore:"twilioToken"        json:"twilioToken"`
	TwilioNumber      string `datastore:"twilioNumber"       json:"twilioNumber"`
}

var globalConfig = &Config{
	Active: false,
}

func getConfig(c *appengine.Context) (*Config, error) {
	// TODO timeout: http://stackoverflow.com/questions/3777367/what-is-a-good-place-to-store-configuration-in-google-appengine-python
	if globalConfig.Active {
		return globalConfig, nil
	} else {
		key := datastore.NewKey(*c, "Config", "TimeVaultConfig", 0, nil)
		if err := datastore.Get(*c, key, globalConfig); err == datastore.ErrNoSuchEntity {
			if _, err := datastore.Put(*c, key, globalConfig); err != nil {
				return globalConfig, err
			}
			return globalConfig, errors.New("No configuration found")
		} else if err != nil {
			return globalConfig, err
		} else if !globalConfig.Active {
			return globalConfig, errors.New("No configuration found")
		} else {
			return globalConfig, nil
		}
	}
}
