package timevault

import (
        "net/http"
        "time"

        "appengine"
        "appengine/datastore"
        "appengine/user"

        "fmt"
)

type TimevaultUser struct {
  ID string
  PhoneNumber string
  Created time.Time
}

func init() {
        http.HandleFunc("/", root)
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
              ID: u.ID,
              PhoneNumber: "555-5555",
              Created: time.Now(),
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

          fmt.Fprintf(w, `Hello, %s `, u)
          fmt.Fprintf(w, `Phone #: %s`, curUser.PhoneNumber)
          return
        } else {
          url, _ := user.LoginURL(c, "/")
          fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
          return
        }
}
