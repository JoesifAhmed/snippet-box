package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/JoesifAhmed/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *applicatoin) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s /n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *applicatoin) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *applicatoin) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *applicatoin) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func (app *applicatoin) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]

	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *applicatoin) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(ContextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
