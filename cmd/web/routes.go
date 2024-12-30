package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/syamimhazmi/snippetbox/ui"
)

// The routes() method returns a servemux containing our application routes
func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamicMiddleware := alice.New(
		app.sessionManager.LoadAndSave,
		noSurf,
		app.authenticate,
	)

	// Guest middleware
	mux.Handle(
		"GET /signup",
		dynamicMiddleware.ThenFunc(app.signup),
	)
	mux.Handle(
		"POST /signup",
		dynamicMiddleware.ThenFunc(app.signupPost),
	)
	mux.Handle(
		"GET /login",
		dynamicMiddleware.ThenFunc(app.login),
	)
	mux.Handle(
		"POST /login",
		dynamicMiddleware.ThenFunc(app.loginPost),
	)
	mux.Handle(
		"GET /{$}",
		dynamicMiddleware.ThenFunc(app.home),
	)
	mux.Handle(
		"GET /snippets/view/{id}",
		dynamicMiddleware.ThenFunc(app.snippetView),
	)

	authMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	mux.Handle(
		"GET /snippets/create",
		authMiddleware.ThenFunc(app.snippetCreate),
	)
	mux.Handle(
		"POST /snippets/store",
		authMiddleware.ThenFunc(app.snippetStore),
	)
	mux.Handle(
		"POST /logout",
		authMiddleware.ThenFunc(app.logout),
	)

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standardMiddleware.Then(mux)
}
