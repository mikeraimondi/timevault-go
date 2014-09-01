package timevault

import (
	"errors"

	"appengine"
	"appengine/datastore"
)

type Config struct {
	Active            bool
	TwilioURL         string
	TwilioMessagePath string
	TwilioSID         string
	TwilioToken       string
	TwilioNumber      string
}

var globalConfig = &Config{
	Active: false,
}

func getConfig(c *appengine.Context) (*Config, error) {
	// TODO timeout: http://stackoverflow.com/questions/3777367/what-is-a-good-place-to-store-configuration-in-google-appengine-python
	if globalConfig.Active {
		return globalConfig, nil
	} else {
		q := datastore.NewQuery("Config").Limit(1)
		var configs []Config
		if _, err := q.GetAll(*c, &configs); err != nil {
			return nil, err
		} else if len(configs) == 0 {
			key := datastore.NewIncompleteKey(*c, "Config", nil)
			if _, err := datastore.Put(*c, key, globalConfig); err != nil {
				return nil, err
			}
			return nil, errors.New("No configuration found")
		} else if !configs[0].Active {
			return nil, errors.New("No configuration found")
		} else {
			globalConfig = &configs[0]
			return globalConfig, nil
		}
	}
}
