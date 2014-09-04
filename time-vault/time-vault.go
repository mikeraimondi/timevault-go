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
	"appengine/taskqueue"
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
		key := datastore.NewKey(c, "TimevaultUser", u.ID, 0, nil)
		curUser := &TimevaultUser{}
		if err := datastore.Get(c, key, curUser); err == datastore.ErrNoSuchEntity {
			newUser := &TimevaultUser{
				// TODO shouldn't need to store the key
				OwnKey:    key,
				Email:     u.Email,
				Username:  u.String(),
				CreatedAt: time.Now(),
			}
			if _, err := datastore.Put(c, key, newUser); err != nil {
				return newUser, err
			}
			return newUser, nil
		} else if err != nil {
			return curUser, err
		} else {
			return curUser, nil
		}
	} else {
		return &TimevaultUser{}, errors.New("Not Logged In")
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
		q := datastore.NewQuery("Pomodoro").Filter("User =", currentUser.OwnKey).Order("-CreatedAt").Limit(1000)
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
		User:      currentUser.OwnKey,
		Duration:  time.Duration(duration) * time.Second,
		CreatedAt: time.Now(),
		Finished:  false,
	}
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Pomodoro", nil)
	if key, err := datastore.Put(c, key, newPomodoro); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// TODO use Delay package?
		t := taskqueue.NewPOSTTask("/endpomodoro", map[string][]string{"key": {key.Encode()}})
		t.Delay = newPomodoro.Duration
		if _, err := taskqueue.Add(c, t, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pom, err := json.Marshal(newPomodoro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(pom), key.Encode())
	}
}

// TODO admin-only auth
func endPomodoro(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c := appengine.NewContext(r)
		config, err := getConfig(&c)
		if err != nil {
			c.Errorf("%v", err)
			return
		}
		var pom Pomodoro
		if key, err := datastore.DecodeKey(r.FormValue("key")); err != nil {
			c.Errorf("%v", err)
			return
		} else {
			if err := datastore.Get(c, key, &pom); err != nil {
				c.Errorf("%v", err)
				return
			}
			if pom.Finished {
				c.Warningf("%v", "Pomodoro already completed")
				return
			}
			pom.Finished = true
			pom.FinishedAt = time.Now()
			if _, err := datastore.Put(c, key, &pom); err != nil {
				c.Errorf("%v", err)
				return
			}
			var u TimevaultUser
			if err := datastore.Get(c, pom.User, &u); err != nil {
				c.Errorf("%v", err)
				return
			}
			if _, err := sendTwilioMessage(config.TwilioNumber, u.PhoneNumber, "Pomodoro complete!", &c); err != nil {
				c.Errorf("%v", err)
				return
			}
			return
		}
	}
}
