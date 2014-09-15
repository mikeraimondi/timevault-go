package timevault

import (
	"encoding/json"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func init() {
	http.HandleFunc("/api/users", withAuth(users))
	http.HandleFunc("/api/pomodoros", withAuth(pomodoros))
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

func pomodoros(w http.ResponseWriter, r *http.Request, u *TimevaultUser) {
	c := appengine.NewContext(r)
	if r.Method != "POST" {
		// Index
		pomodoros, err := u.Pomodoros(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(pomodoros)
		return
	}
	// Create
	duration, err := strconv.ParseInt(r.FormValue("duration"), 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := u.NewPomodoro(&c, duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(p)
}
