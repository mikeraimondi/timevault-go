package timevault

import (
	"appengine/datastore"
	"time"
)

type Pomodoro struct {
	User       *datastore.Key
	Duration   int64
	CreatedAt  time.Time
	Finished   bool
	FinishedAt time.Time
}
