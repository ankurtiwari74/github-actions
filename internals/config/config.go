package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	AppCache      bool
	IsProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}
