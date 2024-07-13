package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Chart and values data
	chart := NewChart()
	values := NewValues()
	templateData := &TemplateData{
		Chart:  chart,
		Values: values,
	}
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// templateCache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
	}

	// flags
	gitlabToken := flag.String("token", "", "Gitlab token for repo")
	flag.Parse()
	if *gitlabToken == "" {
		logger.Error("Missing flag gitlabToken")
		os.Exit(1)
	}

	// Initialize main app struct
	app := NewApplication(
		logger,
		templateData,
		gitlabToken,
		templateCache,
	)

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
