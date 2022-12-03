package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.dasecure.com/internal/models"
	"snippetbox.dasecure.com/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		// ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		// ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		// if err != nil {
		// 	return nil, err
		// }
		// ts, err = ts.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		// 		cache[name] = ts
		// 	}
		// 	return cache,nil
		// }
		// 		files := []string{
		// 			"./ui/html/base.tmpl",
		// 			"./ui/html/partials/nav.tmpl",
		// 			page,
		// 		}

		// 		ts, err := template.ParseFiles(files...)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		cache[name] = ts
	}
	return cache, nil
}
