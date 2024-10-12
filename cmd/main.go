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
	sessionstorage "github.com/fulcrum29/fulcrum/internal/session-storage"
	"github.com/fulcrum29/fulcrum/pkg/templatedata"
	_ "modernc.org/sqlite"
)

func init() {
	// Register TemplateData struct to store in session manager
	gob.Register(templatedata.TemplateData{})
	gob.Register(map[string]interface{}{})
}

func main() {
	// Chart and values data
	templateData := templatedata.NewTemplateData()
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
	sessionstorage.EnsureSessionsTableExists(db)
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 1 * time.Hour

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
