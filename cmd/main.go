package main

import (
  "mntk/internal/handler/http"
  "mntk/internal/pkg/authx"
  "mntk/internal/provider/sqlite/prognosis"
  "mntk/internal/provider/sqlite/question"
  "mntk/internal/provider/sqlite/user"
  "os"
  "os/signal"
  "syscall"

  log "github.com/sirupsen/logrus"
)

const (
  sqliteSource    = "sqlite/dump.db"
  httpHandlerPort = 8080
)

func main() {
  userProvider, err := user.NewProvider(user.Config{
    DataSource: sqliteSource,
  })
  if err != nil {
    log.Fatalf("user.NewProvider: %v", err)
  }

  prognosisProvider, err := prognosis.NewProvider(prognosis.Config{
    DataSource: sqliteSource,
  })
  if err != nil {
    log.Fatalf("prognosis.NewProvider: %v", err)
  }

  questionProvider, err := question.NewProvider(question.Config{
    DataSource: sqliteSource,
  })
  if err != nil {
    log.Fatalf("question.NewProvider: %v", err)
  }

  authProvider := authx.NewProvider(authx.Dependencies{
    User: userProvider,
  })

  httpHandler := http.NewHandler(
    http.Config{
      Port: httpHandlerPort,
    },
    http.Dependencies{
      Providers: http.Providers{
        Auth:      authProvider,
        Question:  questionProvider,
        Prognosis: prognosisProvider,
      },
    },
  )

  go httpHandler.Run()
  wait()
}

func wait() {
  exit := make(chan os.Signal)
  signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
  <-exit
}
