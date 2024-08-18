package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "modernc.org/sqlite"
)

func init() {
	gob.Register(TemplateData{})
}

func main() {
	// Chart and values data
	templateData := NewTemplateData()
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// html templates cache
	htmlTemplateCache, err := newHtmlTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// db connection
	db, err := sql.Open("sqlite", "fulcrum.db")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// session manager
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// flags
	gitlabToken := flag.String("token", "", "Gitlab token for repo")
	flag.Parse()
	if *gitlabToken == "" {
		flag.PrintDefaults()
		logger.Error("Missing flag gitlabToken")
		os.Exit(1)
	}

	// Initialize main app
	app := &application{
		logger:            logger,
		templateData:      templateData,
		gitlabToken:       gitlabToken,
		htmlTemplateCache: htmlTemplateCache,
		sessionManager:    sessionManager,
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	// Start main handler (server)
	logger.Info("Starting server")

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
