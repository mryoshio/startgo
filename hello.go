package guestbook

import (
  "html/template"
  "net/http"
  "time"

  "appengine"
  "appengine/datastore"
  "appengine/user"
)

type Greeting struct {
  Author string
  Content string
  Date time.Time
}

var (
  gbTemplate = template.Must(template.ParseFiles("view/guestbook.html"))
)

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/sign", sign)
}

func guestbookKey(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "GuestBook", "default_guestbook", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
  greetings := make([]Greeting, 0, 10)

  if _, err := q.GetAll(c, &greetings); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  if err := gbTemplate.Execute(w, greetings); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func sign(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  g := Greeting {
    Content: r.FormValue("content"),
    Date:    time.Now(),
  }

  if u := user.Current(c); u != nil {
    g.Author = u.String()
  }

  key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
  _, err := datastore.Put(c, key, &g)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, "/", http.StatusFound)
}
