package main

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/fulcrum29/fulcrum/pkg/templatedata"
)

type application struct {
	logger            *slog.Logger
	templateData      *templatedata.TemplateData
	gitlabToken       *string
	htmlTemplateCache map[string]*template.Template
	sessionManager    *scs.SessionManager
}

// func NewApplication(
// 	l *slog.Logger,
// 	td *TemplateData,
// 	gt *string,
// 	tc map[string]*template.Template,
// ) *application {
// 	return &application{
// 		logger:            l,
// 		templateData:      td,
// 		gitlabToken:       gt,
// 		htmlTemplateCache: tc,
// 	}
// }
