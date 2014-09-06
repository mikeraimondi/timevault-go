package timevault

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"appengine/taskqueue"
	"appengine/user"
)

func init() {
	http.HandleFunc("/users", withAuth(users))
	http.HandleFunc("/pomodoros", withAuth(index))
}

func withAuth(handler func(http.ResponseWriter, *http.Request, *TimevaultUser)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u, err := authenticate(w, r); err != nil {
			if err.Error() == "Not Logged In" {
				c := appengine.NewContext(r)
				url, _ := user.LoginURL(c, r.URL.String())
				http.Redirect(w, r, url, 302)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			handler(w, r, u)
		}
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) (*TimevaultUser, error) {
	c := appengine.NewContext(r)
	u, err := NewTimevaultUser(&c)
	// TODO cleanup
	if err != nil {
		return u, err
	}
	if err := datastore.Get(c, u.CachedKey, u); err == datastore.ErrNoSuchEntity {
		if _, err := datastore.Put(c, u.CachedKey, u); err != nil {
			return u, err
		}
		return u, nil
	} else if err != nil {
		return u, err
	} else {
		return u, nil
	}
}

func users(w http.ResponseWriter, r *http.Request, u *TimevaultUser) {
	//TODO make RESTful
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(u)
	return
}

func index(w http.ResponseWriter, r *http.Request, u *TimevaultUser) {
	if r.Method != "POST" {
		c := appengine.NewContext(r)
		// TODO Paginate
		q := datastore.NewQuery("Pomodoro").Ancestor(u.CachedKey).Order("-createdAt").Limit(1000)
		pomodoros := make([]Pomodoro, 0, 1000)
		if _, err := q.GetAll(c, &pomodoros); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(pomodoros)
		return
	}
	duration, err := strconv.ParseInt(r.FormValue("duration"), 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := appengine.NewContext(r)
	newPomodoro := &Pomodoro{
		Duration:  time.Duration(duration) * time.Second,
		CreatedAt: time.Now(),
		Finished:  false,
	}
	key := datastore.NewIncompleteKey(c, "Pomodoro", u.CachedKey)
	key, err = datastore.Put(c, key, newPomodoro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := endPomodoro.Task(*key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Delay = newPomodoro.Duration
	if _, err := taskqueue.Add(c, t, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(newPomodoro)
}

var endPomodoro = delay.Func("endPomodoro", func(c appengine.Context, k datastore.Key) error {
	config, err := getConfig(&c)
	if err != nil {
		return err
	}
	var pom Pomodoro
	if err := datastore.Get(c, &k, &pom); err != nil {
		return err
	}
	if pom.Finished {
		c.Warningf("%v", "Pomodoro already completed")
		return nil
	}
	pom.Finished = true
	pom.FinishedAt = time.Now()
	if _, err := datastore.Put(c, &k, &pom); err != nil {
		return err
	}
	var u TimevaultUser
	if err := datastore.Get(c, k.Parent(), &u); err != nil {
		return err
	}
	if _, err := sendTwilioMessage(config.TwilioNumber, u.PhoneNumber, "Pomodoro complete!", &c); err != nil {
		return err
	}
	return nil
})
