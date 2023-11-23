package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *applicatoin) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.Then(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddleware.Then(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamicMiddleware.Then(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddleware.Then(http.HandlerFunc(app.showSnippet)))

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileserver))

	return standardMiddleware.Then(mux)
}
