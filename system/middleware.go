package system

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/golang/glog"
	"github.com/gorilla/sessions"
	"github.com/haruyama/golang-goji-sample/models"
	"github.com/zenazn/goji/web"
)

// Makes sure templates are stored in the context
func (application *Application) ApplyTemplates(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["Template"] = application.Template
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Makes sure controllers can have access to session
func (application *Application) ApplySessions(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := application.Store.Get(r, "session")
		c.Env["Session"] = session
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) ApplyDbMap(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["DbMap"] = application.DbMap
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := c.Env["Session"].(*sessions.Session)
		if userId, ok := session.Values["UserId"]; ok {
			dbMap := c.Env["DbMap"].(*gorp.DbMap)

			user, err := dbMap.Get(models.User{}, userId)
			if err != nil {
				glog.Warningf("Auth error: %v", err)
				c.Env["User"] = nil
			} else {
				c.Env["User"] = user
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) ApplyIsXhr(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			c.Env["IsXhr"] = true
		} else {
			c.Env["IsXhr"] = false
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
