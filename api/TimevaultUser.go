package api

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

func (u *TimevaultUser) Pomodoros(c *appengine.Context) ([]Pomodoro, error) {
	q := datastore.NewQuery("Pomodoro").Ancestor(u.CachedKey).Order("-createdAt").Limit(1000)
	var pomodoros []Pomodoro
	_, err := q.GetAll(*c, &pomodoros)
	return pomodoros, err
}

func (u *TimevaultUser) NewPomodoro(c *appengine.Context, duration int64) (*Pomodoro, error) {
	p := &Pomodoro{
		Duration:  time.Duration(duration) * time.Second,
		CreatedAt: time.Now(),
		Finished:  false,
	}
	key := datastore.NewIncompleteKey(*c, "Pomodoro", u.CachedKey)
	key, err := datastore.Put(*c, key, p)
	if err != nil {
		return p, err
	}
	if err := p.Finish(c, key); err != nil {
		return p, err
	}
	return p, nil
}
