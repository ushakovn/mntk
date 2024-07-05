package user

import (
  "context"
  "database/sql"
  "fmt"
  "mntk/internal/models/user"
  "mntk/internal/pkg/sqlx"

  sq "github.com/Masterminds/squirrel"

  _ "github.com/glebarez/go-sqlite"
)

const tableName = "users"

var tableColumns = sqlx.Columns(user.User{})

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

func (s *Provider) ChangePassword(ctx context.Context, params user.ChangePasswordParams) (*user.User, error) {
  suffix := sqlx.ReturningSuffix(tableColumns)

  builder := sq.Update(tableName).
    Set("password", params.Password).
    Where(sq.Eq{"name": params.Name}).
    Suffix(suffix)

  return sqlx.Get[user.User](ctx, s.conn, builder)
}

func (s *Provider) GetByName(ctx context.Context, name string) (*user.User, error) {
  builder := sq.Select(tableColumns).
    From(tableName).
    Where(sq.Eq{"name": name})

  return sqlx.Get[user.User](ctx, s.conn, builder)
}
