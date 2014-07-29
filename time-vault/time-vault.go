package timevault

import (
	"net/http"
	"time"
	"encoding/json"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"fmt"
)

type TimevaultUser struct {
	ID          string
	Email       string
	Username    string
	PhoneNumber string
	Created     time.Time
}

type Pomodoro struct {
	User string
	Duration int64
	Created time.Time
}

func (u *TimevaultUser) String() string {
	return u.Username
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/pomodoros", index)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if u := user.Current(c); u != nil {
		q := datastore.NewQuery("TimevaultUser").Filter("ID =", u.ID).Limit(1)
		users := make([]TimevaultUser, 0, 1)
		if _, err := q.GetAll(c, &users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		curUser := &TimevaultUser{}
		if len(users) == 0 {
			newUser := &TimevaultUser{
				ID:          u.ID,
				Email:       u.Email,
				Username:    u.String(),
				PhoneNumber: "555-5555",
				Created:     time.Now(),
			}
			key := datastore.NewIncompleteKey(c, "TimevaultUser", nil)
			if _, err := datastore.Put(c, key, newUser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			curUser = newUser
		} else {
			curUser = &users[0]
		}
		us, err := json.Marshal(curUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(us))
	} else {
		url, _ := user.LoginURL(c, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if u := user.Current(c); u != nil {
		q := datastore.NewQuery("TimevaultUser").Filter("ID =", u.ID).Limit(1)
		users := make([]TimevaultUser, 0, 1)
		if _, err := q.GetAll(c, &users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		curUser := &TimevaultUser{}
		if len(users) == 0 {
			newUser := &TimevaultUser{
				ID:          u.ID,
				Email:       u.Email,
				Username:    u.String(),
				PhoneNumber: "555-5555",
				Created:     time.Now(),
			}
			key := datastore.NewIncompleteKey(c, "TimevaultUser", nil)
			if _, err := datastore.Put(c, key, newUser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			curUser = newUser
		} else {
			curUser = &users[0]
		}
		if r.Method != "POST" {
			// index
			return
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		duration, err := strconv.ParseInt(r.Form["duration"][0], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newPomodoro := &Pomodoro{
			User: curUser.ID,
			Duration: duration,
			Created: time.Now(),
		}
		key := datastore.NewIncompleteKey(c, "Pomodoro", nil)
		if _, err := datastore.Put(c, key, newPomodoro); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pom, err := json.Marshal(newPomodoro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(pom))
	} else {
		url, _ := user.LoginURL(c, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
	}
}
