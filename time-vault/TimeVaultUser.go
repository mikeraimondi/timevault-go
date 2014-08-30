package timevault

import "time"

type TimevaultUser struct {
  ID          string
  Email       string
  Username    string
  PhoneNumber string
  CreatedAt   time.Time
}

func (u *TimevaultUser) String() string {
  return u.Username
}
