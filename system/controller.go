package system

import (
	"bytes"
	"html/template"

	"gopkg.in/gorp.v1"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
)

type Controller struct {
}

func (controller *Controller) GetSession(c web.C) *sessions.Session {
	return c.Env["Session"].(*sessions.Session)
}

func (controller *Controller) GetTemplate(c web.C) *template.Template {
	return c.Env["Template"].(*template.Template)
}

func (controller *Controller) GetDbMap(c web.C) *gorp.DbMap {
	return c.Env["DbMap"].(*gorp.DbMap)
}

func (controller *Controller) IsXhr(c web.C) bool {
	return c.Env["IsXhr"].(bool)
}

func (controller *Controller) Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
