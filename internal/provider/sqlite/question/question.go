package question

import (
  "context"
  "database/sql"
  "fmt"
  "mntk/internal/models/question"
  "mntk/internal/pkg/sqlx"

  sq "github.com/Masterminds/squirrel"

  _ "github.com/glebarez/go-sqlite"
)

const tableName = "questions"

var tableColumns = sqlx.Columns(question.Question{})

type Provider struct {
  conn *sql.DB
}

type Config struct {
  DataSource string
}

func NewProvider(config Config) (*Provider, error) {
  conn, err := sql.Open("sqlite", config.DataSource)
  if err != nil {
    return nil, fmt.Errorf("sql.Open: %w", err)
  }

  if err = conn.Ping(); err != nil {
    return nil, fmt.Errorf("conn.Ping: %w", err)
  }

  return &Provider{
    conn: conn,
  }, nil
}

func (s *Provider) Create(ctx context.Context, params question.CreateParams) (*question.Question, error) {
  suffix := sqlx.ReturningSuffix(tableColumns)

  fields := map[string]any{
    "prognosis_id": params.PrognosisID,
    "label":        params.Label,
    "answer":       params.Answer,
    "score":        params.Score,
  }

  builder := sq.Insert(tableName).
    SetMap(fields).
    Suffix(suffix)

  return sqlx.Get[question.Question](ctx, s.conn, builder)
}
