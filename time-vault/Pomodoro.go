package timevault

import (
	"time"
)

type Pomodoro struct {
	Duration   time.Duration `datastore:"duration"    json:"duration"`
	CreatedAt  time.Time     `datastore:"createdAt"   json:"createdAt"`
	Finished   bool          `datastore:"finished"    json:"finished"`
	FinishedAt time.Time     `datastore:"finishedAt"  json:"finishedAt"`
}
