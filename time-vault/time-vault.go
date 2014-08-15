package timevault

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func (u *TimevaultUser) Error() string {
	if u.ID == "" {
		return "Not Logged In"
	}
	return ""
}

func init() {
	http.HandleFunc("/", withAuth(root))
	http.HandleFunc("/pomodoros", withAuth(index))
}

func withAuth(handler func(http.ResponseWriter, *http.Request, *TimevaultUser)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if curUser, err := authenticate(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if curUser.Error() == "Not Logged In" {
				return
			}
			handler(w, r, curUser)
		}
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) (*TimevaultUser, error) {
	c := appengine.NewContext(r)
	if u := user.Current(c); u != nil {
		q := datastore.NewQuery("TimevaultUser").Filter("ID =", u.ID).Limit(1)
		users := make([]TimevaultUser, 0, 1)
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
		url, _ := user.LoginURL(c, "/")
		http.Redirect(w, r, url, 302)
		return &TimevaultUser{}, nil
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
