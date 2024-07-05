package http

import (
  "context"
  "mntk/internal/models/prognosis"
  "mntk/internal/models/question"
  "net/http"
)

type Dependencies struct {
  Providers Providers
}

type Providers struct {
  Auth      AuthProvider
  Question  QuestionProvider
  Prognosis PrognosisProvider
}

type PrognosisProvider interface {
  List(ctx context.Context, params prognosis.ListParams) ([]*prognosis.Prognosis, error)
  Create(ctx context.Context, params prognosis.CreateParams) (*prognosis.Prognosis, error)
}

type QuestionProvider interface {
  Create(ctx context.Context, params question.CreateParams) (*question.Question, error)
}

type AuthProvider interface {
  Middleware(handler http.HandlerFunc) http.HandlerFunc
  ChangePassword(ctx context.Context, oldPassword, newPassword string) error
}
