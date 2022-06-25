package render

import (
	"bookings/internals/config"
	"bookings/internals/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var function = template.FuncMap{}
var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// cache, err := CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var cache map[string]*template.Template
	if app.AppCache {
		cache = app.TemplateCache

	} else {
		cache, _ = CreateTemplateCache()
		fmt.Println(cache)
	}

	t, ok := cache[tmpl]
	if !ok {
		log.Fatal("Could not get template cache")
	}

	buf := new(bytes.Buffer)
	td = AddTemplateData(td, r)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error rendering the template", err)
	}

}
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	fmt.Println("-----Pages------", pages)
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(function).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		fmt.Println("-----Page------", page)
		base, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(base) > 0 {
			fmt.Println("-----Base------", base)
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
			fmt.Println("-----ts------", ts)
		}
		myCache[name] = ts
	}
	return myCache, nil
}
