package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Razzle131/merchStore/internal/db"
	"github.com/Razzle131/merchStore/internal/handler"
	"github.com/Razzle131/merchStore/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func main() {
	godotenv.Load()

	cfg := handler.Config{
		Port: os.Getenv("SERVER_ADDRESS"),
		DSN:  os.Getenv("POSTGRES_CONN"),
	}

	// init logger
	logger.SetupLogging(slog.LevelDebug)

	slog.Debug("Debugging info enabled")

	slog.Info("Starting db", slog.String("DSN", os.Getenv("POSTGRES_CONN")))

	dbConn, err := db.New(context.Background(), cfg.DSN)
	if err != nil {
		log.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbConn.Close()

	myServer := handler.NewServer(dbConn)

	r := http.NewServeMux()

	r.HandleFunc("POST /api/auth", myServer.PostApiAuth)
	r.HandleFunc("GET /api/buy/{item}", func(w http.ResponseWriter, r *http.Request) {
		item := strings.Split(r.URL.Path, "/")[3]
		myServer.GetApiBuyItem(w, r, item)
	})
	r.HandleFunc("GET /api/info", myServer.GetApiInfo)
	r.HandleFunc("POST /api/sendCoin", myServer.PostApiSendCoin)

	srv := &http.Server{
		Addr:    "localhost:" + cfg.Port,
		Handler: r,
	}
	slog.Info("Starting server on address " + srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
