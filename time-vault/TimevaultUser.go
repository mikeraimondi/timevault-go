package timevault

import (
	"errors"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type TimevaultUser struct {
	Email       string         `datastore:"email"       json:"email"`
	Username    string         `datastore:"username"    json:"username"`
	PhoneNumber string         `datastore:"phoneNumber" json:"phoneNumber"`
	CreatedAt   time.Time      `datastore:"createdAt"   json:"createdAt"`
	CachedKey   *datastore.Key `datastore:"-"           json:"-"`
}

func NewTimevaultUser(c *appengine.Context) (*TimevaultUser, error) {
	if u := user.Current(*c); u != nil {
		return &TimevaultUser{
			CachedKey: datastore.NewKey(*c, "TimevaultUser", u.ID, 0, nil),
			Email:     u.Email,
			Username:  u.String(),
			CreatedAt: time.Now(),
		}, nil
	}
	return &TimevaultUser{}, errors.New("Not Logged In")
}

func (u *TimevaultUser) String() string {
	return u.Username
}
