package timevault

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type TimevaultUser struct {
	Email           string    `datastore:"email"       json:"email"`
	Username        string    `datastore:"username"    json:"username"`
	PhoneNumber     string    `datastore:"phoneNumber" json:"phoneNumber"`
	CreatedAt       time.Time `datastore:"createdAt"   json:"createdAt"`
	GoogleAccountID string    `datastore:"-"           json:"-"`
}

func (u *TimevaultUser) String() string {
	return u.Username
}

func (u *TimevaultUser) Key(c *appengine.Context) *datastore.Key {
	return datastore.NewKey(*c, "TimevaultUser", u.GoogleAccountID, 0, nil)
}
