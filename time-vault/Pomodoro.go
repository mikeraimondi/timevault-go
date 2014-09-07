package timevault

import (
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
)

type Pomodoro struct {
	Duration   time.Duration `datastore:"duration"    json:"duration"`
	CreatedAt  time.Time     `datastore:"createdAt"   json:"createdAt"`
	Finished   bool          `datastore:"finished"    json:"finished"`
	FinishedAt time.Time     `datastore:"finishedAt"  json:"finishedAt"`
}

func (p *Pomodoro) Finish(c *appengine.Context, k *datastore.Key) error {
	t, err := endPomodoro.Task(*k)
	if err != nil {
		return err
	}
	t.Delay = p.Duration
	if _, err := taskqueue.Add(*c, t, ""); err != nil {
		return err
	}
	return nil
}

var endPomodoro = delay.Func("endPomodoro", func(c appengine.Context, k datastore.Key) error {
	config, err := getConfig(&c)
	if err != nil {
		return err
	}
	var pom Pomodoro
	if err := datastore.Get(c, &k, &pom); err != nil {
		return err
	}
	if pom.Finished {
		c.Warningf("%v", "Pomodoro already completed")
		return nil
	}
	pom.Finished = true
	pom.FinishedAt = time.Now()
	if _, err := datastore.Put(c, &k, &pom); err != nil {
		return err
	}
	var u TimevaultUser
	if err := datastore.Get(c, k.Parent(), &u); err != nil {
		return err
	}
	if _, err := sendTwilioMessage(config.TwilioNumber, u.PhoneNumber, "Pomodoro complete!", &c); err != nil {
		return err
	}
	return nil
})
