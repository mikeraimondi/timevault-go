package timevault

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"errors"

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
	User     string
	Duration int64
	Created  time.Time
}

func (u *TimevaultUser) String() string {
	return u.Username
}

func init() {
	http.HandleFunc("/", withAuth(root))
	http.HandleFunc("/pomodoros", withAuth(index))
}

func withAuth(handler func(http.ResponseWriter, *http.Request, *TimevaultUser)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if curUser, err := authenticate(w, r); err != nil {
			if err.Error() == "Not Logged In" {
				c := appengine.NewContext(r)
				url, _ := user.LoginURL(c, "/")
				http.Redirect(w, r, url, 302)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			handler(w, r, curUser)
		}
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) (*TimevaultUser, error) {
	c := appengine.NewContext(r)
	if u := user.Current(c); u != nil {
		q := datastore.NewQuery("TimevaultUser").Filter("ID =", u.ID).Limit(1)
		var users []TimevaultUser
		if _, err := q.GetAll(c, &users); err != nil {
			return nil, err
		}
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
				return nil, err
			}
			return newUser, nil
		} else {
			return &users[0], nil
		}
	} else {
		return nil, errors.New("Not Logged In")
	}
}

func root(w http.ResponseWriter, r *http.Request, curUser *TimevaultUser) {
	if us, err := json.Marshal(curUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(us))
		return
	}
}

func index(w http.ResponseWriter, r *http.Request, curUser *TimevaultUser) {
	if r.Method != "POST" {
		c := appengine.NewContext(r)
		// Paginate
		q := datastore.NewQuery("Pomodoro").Filter("User =", curUser.ID).Order("-Created").Limit(1000)
		pomodoros := make([]Pomodoro, 0, 1000)
		if _, err := q.GetAll(c, &pomodoros); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if poms, err := json.Marshal(pomodoros); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(poms))
			return
		}
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
		User:     curUser.ID,
		Duration: duration,
		Created:  time.Now(),
	}
	c := appengine.NewContext(r)
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
}
