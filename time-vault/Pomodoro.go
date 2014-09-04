package timevault

import (
	"appengine/datastore"
	"time"
)

type Pomodoro struct {
	User       *datastore.Key
	Duration   time.Duration
	CreatedAt  time.Time
	Finished   bool
	FinishedAt time.Time
}
