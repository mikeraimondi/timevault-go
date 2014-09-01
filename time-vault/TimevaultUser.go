package timevault

import (
	"appengine/datastore"
	"time"
)

type TimevaultUser struct {
	OwnKey      *datastore.Key
	Email       string
	Username    string
	PhoneNumber string
	CreatedAt   time.Time
}

func (u *TimevaultUser) String() string {
	return u.Username
}
