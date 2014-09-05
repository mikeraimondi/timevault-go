package timevault

import (
	"appengine/datastore"
	"time"
)

type Pomodoro struct {
	User       *datastore.Key `datastore:"user"        json:"user"`
	Duration   time.Duration  `datastore:"duration"    json:"duration"`
	CreatedAt  time.Time      `datastore:"createdAt"   json:"createdAt"`
	Finished   bool           `datastore:"finished"    json:"finished"`
	FinishedAt time.Time      `datastore:"finishedAt"  json:"finishedAt"`
}
