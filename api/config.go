package api

import (
	"errors"
	"time"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/sessions"
)

type Config struct {
	Active               bool                  `datastore:"active"               json:"active"`
	GplusApplicationName string                `datastore:"gplusApplicationName" json:"gplusApplicationName"`
	GplusClientID        string                `datastore:"gplusClientID"        json:"gplusClientID"`
	GplusClientSecret    string                `datastore:"gplusClientSecret"    json:"gplusClientSecret"`
	GplusRedirectURL     string                `datastore:"gplusRedirectURL"     json:"gplusRedirectURL"`
	SessionSecret        string                `datastore:"sessionSecret"        json:"sessionSecret"`
	TwilioURL            string                `datastore:"twilioURL"            json:"twilioURL"`
	TwilioMessagePath    string                `datastore:"twilioMessagePath"    json:"twilioMessagePath"`
	TwilioSID            string                `datastore:"twilioSID"            json:"twilioSID"`
	TwilioToken          string                `datastore:"twilioToken"          json:"twilioToken"`
	TwilioNumber         string                `datastore:"twilioNumber"         json:"twilioNumber"`
	Store                *sessions.CookieStore `datastore:"-"`
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
			// TODO memoize/cache
			globalConfig.Store = sessions.NewCookieStore([]byte(globalConfig.SessionSecret))
			globalConfig.Store.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   int((time.Hour * 24 * 7) / time.Second),
				HttpOnly: true,
			}

			return globalConfig, nil
		}
	}
}
