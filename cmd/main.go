package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	values := NewValues()
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// flags
	gitlabToken := flag.String("token", "", "Gitlab token for repo")
	flag.Parse()
	if *gitlabToken == "" {
		logger.Error("Missing flag gitlabToken")
		os.Exit(1)
	}

	app := NewApplication(
		logger,
		values,
		gitlabToken)

	// Start main handler (server)
	logger.Info("Starting server")
	err := http.ListenAndServe(":8080", app.routes())

	logger.Error(err.Error())
	os.Exit(1)

}
