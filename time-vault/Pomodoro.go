package timevault

import "time"

type Pomodoro struct {
  User       string
  Duration   int64
  CreatedAt  time.Time
  Finished   bool
  FinishedAt time.Time
}
