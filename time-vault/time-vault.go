package timevault

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", withAuth(root))
	http.HandleFunc("/pomodoros", withAuth(index))
	http.HandleFunc("/endpomodoro", endPomodoro)
}

func withAuth(handler func(http.ResponseWriter, *http.Request, *TimevaultUser)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if currentUser, err := authenticate(w, r); err != nil {
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
			handler(w, r, currentUser)
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
				ID:        u.ID,
				Email:     u.Email,
				Username:  u.String(),
				CreatedAt: time.Now(),
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

func root(w http.ResponseWriter, r *http.Request, currentUser *TimevaultUser) {
	if us, err := json.Marshal(currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(us))
		return
	}
}

func index(w http.ResponseWriter, r *http.Request, currentUser *TimevaultUser) {
	if r.Method != "POST" {
		c := appengine.NewContext(r)
		// TODO Paginate
		q := datastore.NewQuery("Pomodoro").Filter("User =", currentUser.ID).Order("-CreatedAt").Limit(1000)
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
	duration, err := strconv.ParseInt(r.FormValue("duration"), 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newPomodoro := &Pomodoro{
		User:      currentUser.ID,
		Duration:  duration,
		CreatedAt: time.Now(),
		Finished:  false,
	}
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Pomodoro", nil)
	if key, err := datastore.Put(c, key, newPomodoro); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		pom, err := json.Marshal(newPomodoro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// TODO Enqueue worker with key
		fmt.Fprint(w, string(pom), key.Encode())
	}
}

func endPomodoro(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c := appengine.NewContext(r)
		var pom Pomodoro
		if key, err := datastore.DecodeKey(r.FormValue("key")); err != nil {
			c.Errorf("%v", err)
			return
		} else {
			if err := datastore.Get(c, key, &pom); err != nil {
				c.Errorf("%v", err)
				return
			}
			pom.Finished = true
			pom.FinishedAt = time.Now()
			if _, err := datastore.Put(c, key, &pom); err != nil {
				c.Errorf("%v", err)
				return
			}
			config, err := getConfig(&c)
			if err != nil {
				c.Errorf("%v", err)
				return
			}
			// TODO hardcoded values
			if _, err := sendTwilioMessage(config.TwilioNumber, "+16787612326", "Pomodoro complete!", &c); err != nil {
				c.Errorf("%v", err)
				return
			}
			return
		}
	}
}
