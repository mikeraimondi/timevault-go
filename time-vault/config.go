package timevault

import (
	"errors"

	"appengine"
	"appengine/datastore"
)

type Config struct {
	TwilioURL         string
	TwilioMessagePath string
	TwilioSID         string
	TwilioToken       string
	Active            bool
}

var config = &Config{
	Active: false,
}

func getConfig(c appengine.Context) (*Config, error) {
	if config.Active {
		return config, nil
	} else {
		q := datastore.NewQuery("Config").Limit(1)
		var configs []Config
		if _, err := q.GetAll(c, &configs); err != nil {
			return nil, err
		} else if len(configs) == 0 {
			key := datastore.NewIncompleteKey(c, "Config", nil)
			if _, err := datastore.Put(c, key, config); err != nil {
				return nil, err
			}
			return nil, errors.New("No configuration found")
		} else if !configs[0].Active {
			return nil, errors.New("No configuration found")
		} else {
			config = &configs[0]
			return config, nil
		}
	}
}
