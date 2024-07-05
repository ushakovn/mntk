package prognosis

import (
  "context"
  "database/sql"
  "fmt"
  "mntk/internal/models/prognosis"
  "mntk/internal/pkg/sqlx"

  sq "github.com/Masterminds/squirrel"

  _ "github.com/glebarez/go-sqlite"
)

const tableName = "prognosis"

var tableColumns = sqlx.Columns(prognosis.Prognosis{})

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

func (s *Provider) List(ctx context.Context, params prognosis.ListParams) ([]*prognosis.Prognosis, error) {
  builder := sq.Select(tableColumns).
    From(tableName)

  pagination := params.Pagination

  if pagination.IsValid() {
    builder = builder.
      Offset(pagination.Offset).
      Limit(pagination.Limit)
  }

  sort := params.Sort

  if sort.IsValid() {
    builder = builder.
      OrderBy(sort.OrderBy())
  }

  return sqlx.Select[prognosis.Prognosis](ctx, s.conn, builder)
}

func (s *Provider) Create(ctx context.Context, params prognosis.CreateParams) (*prognosis.Prognosis, error) {
  suffix := sqlx.ReturningSuffix(tableColumns)

  fields := map[string]any{
    "score":         params.Score,
    "result":        params.Result,
    "patient_name":  params.PatientName,
    "patient_birth": params.PatientBirth,
  }

  builder := sq.Insert(tableName).
    SetMap(fields).
    Suffix(suffix)

  return sqlx.Get[prognosis.Prognosis](ctx, s.conn, builder)
}
