package system

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/go-utils/uslice"
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

func isValidToken(a, b string) bool {
	x := []byte(a)
	y := []byte(b)
	if len(x) != len(y) {
		return false
	}
	return subtle.ConstantTimeCompare(x, y) == 1
}

var csrfProtectionMethodForNoXhr = []string{"POST", "PUT", "DELETE"}

func isCsrfProtectionMethodForNoXhr(method string) bool {
	return uslice.StrHas(csrfProtectionMethodForNoXhr, strings.ToUpper(method))
}

func (application *Application) ApplyCsrfProtection(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := c.Env["Session"].(*sessions.Session)
		csrfProtection := application.CsrfProtection
		if _, ok := session.Values["CsrfToken"]; !ok {
			hash := sha256.New()
			buffer := make([]byte, 32)
			_, err := rand.Read(buffer)
			if err != nil {
				glog.Fatalf("crypt/rand.Read failed: %s", err)
			}
			hash.Write(buffer)
			session.Values["CsrfToken"] = fmt.Sprintf("%x", hash.Sum(nil))
			if err = session.Save(r, w); err != nil {
				glog.Fatal("session.Save() failed")
			}
		}
		c.Env["CsrfKey"] = csrfProtection.Key
		c.Env["CsrfToken"] = session.Values["CsrfToken"]
		csrfToken := c.Env["CsrfToken"].(string)

		if c.Env["IsXhr"].(bool) {
			if !isValidToken(csrfToken, r.Header.Get(csrfProtection.Header)) {
				http.Error(w, "Invalid Csrf Header", http.StatusBadRequest)
				return
			}
		} else {
			if isCsrfProtectionMethodForNoXhr(r.Method) {
				if !isValidToken(csrfToken, r.PostFormValue(csrfProtection.Key)) {
					http.Error(w, "Invalid Csrf Token", http.StatusBadRequest)
					return
				}
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:   csrfProtection.Cookie,
			Value:  csrfToken,
			Secure: csrfProtection.Secure,
		})
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
