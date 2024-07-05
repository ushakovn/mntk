package main

import (
  "mntk/internal/handler/http"
  "mntk/internal/pkg/authx"
  "mntk/internal/storage/sqlite"
  "os"
  "os/signal"
  "syscall"

  log "github.com/sirupsen/logrus"
)

func main() {
  storage, err := sqlite.NewStorage(sqlite.Config{
    DataSource: "sqlite/dump.db",
  })
  if err != nil {
    log.Fatalf("sqlite.NewStorage: %v", err)
  }

  authHandler := authx.NewHandler(authx.Config{
    Storage: storage,
  })

  handler := http.NewHandler(http.Config{
    Port:    8080,
    Storage: storage,
    Auth:    authHandler,
  })

  go handler.Run()

  exit := make(chan os.Signal)
  signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
  <-exit
}
