package main

import (
	"html/template"
	"path/filepath"

	"github.com/syamimhazmi/snippetbox/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// init a new map to act as a cache
	cache := map[string]*template.Template{}

	// use filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". this will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl, ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract file from full filepath
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
