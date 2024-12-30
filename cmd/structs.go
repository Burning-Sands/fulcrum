package main

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/fulcrum29/fulcrum/pkg/templatedata"
	"github.com/google/go-github/v66/github"
)

type application struct {
	logger            *slog.Logger
	templateData      *templatedata.TemplateData
	githubToken       *string
	htmlTemplateCache map[string]*template.Template
	sessionManager    *scs.SessionManager
	githubClient      *github.Client
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
