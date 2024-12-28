package main

import "net/http"

// The routes() method returns a servemux containing our application routes
func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippets/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippets/create", app.snippetCreate)
	mux.HandleFunc("POST /snippets/store", app.snippetStore)

	return app.recoverPanic(
		app.logRequest(
			commonHeaders(mux),
		),
	)
}
