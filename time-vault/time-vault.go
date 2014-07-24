package timevault

import (
        "net/http"

        "appengine"
        "appengine/user"

        "fmt"
)

func init() {
        http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        if u := user.Current(c); u != nil {
          fmt.Fprintf(w, `Hello, %s`, u)
          return
        } else {
          url, _ := user.LoginURL(c, "/")
          fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
          return
        }
}
