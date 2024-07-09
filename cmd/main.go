package main

import (
	"flag"
	"fmt"
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
	fmt.Println(app.templateCache)

	// Start main handler (server)
	logger.Info("Starting server")

	err = http.ListenAndServe(":8080", app.routes())
	logger.Error(err.Error())
	os.Exit(1)

}
