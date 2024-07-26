package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/ErickMachado/snippetbox/internal/models"
	"github.com/ErickMachado/snippetbox/ui"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func brazillianDate(t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Time{}, err
	}

	locTime := t.In(loc)

	return locTime, nil
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"brazillianDate": brazillianDate,
}

type templateData struct {
	Snippet         models.Snippet
	Snippets        []models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
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

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
