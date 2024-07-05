package sqlite

import (
  "context"
  "database/sql"
  "fmt"
  "mntk/internal/models"
  "reflect"
  "strings"

  sq "github.com/Masterminds/squirrel"
  "github.com/georgysavva/scany/v2/sqlscan"
  _ "github.com/glebarez/go-sqlite"
)

const (
  prognosisTableName = "prognosis"
  questionTableName  = "questions"
  usersTableName     = "users"
)

var (
  prognosisTableColumns = tableColumns(models.Prognosis{})
  questionTableColumns  = tableColumns(models.Question{})
  usersTableColumns     = tableColumns(models.User{})
)

type Storage struct {
  conn *sql.DB
}

type Config struct {
  DataSource string
}

func NewStorage(config Config) (*Storage, error) {
  conn, err := sql.Open("sqlite", config.DataSource)
  if err != nil {
    return nil, fmt.Errorf("sql.Open: %w", err)
  }

  if err = conn.Ping(); err != nil {
    return nil, fmt.Errorf("conn.Ping: %w", err)
  }

  return &Storage{
    conn: conn,
  }, nil
}

func (s *Storage) ListPrognosis(ctx context.Context, params models.ListPrognosisParams) ([]*models.Prognosis, error) {
  builder := sq.Select(prognosisTableColumns).
    From(prognosisTableName)

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

  return selectx[models.Prognosis](ctx, s.conn, builder)
}

func (s *Storage) CreatePrognosis(ctx context.Context, params models.CreatePrognosisParams) (*models.Prognosis, error) {
  suffix := returningSuffix(prognosisTableColumns)

  fields := map[string]any{
    "score":  params.Score,
    "result": params.Result,
  }

  if params.PatientName != nil {
    fields["patient_name"] = params.PatientName
  }

  builder := sq.Insert(prognosisTableName).
    SetMap(fields).
    Suffix(suffix)

  return getx[models.Prognosis](ctx, s.conn, builder)
}

func (s *Storage) CreateQuestion(ctx context.Context, params models.CreateQuestionParams) (*models.Question, error) {
  suffix := returningSuffix(questionTableColumns)

  fields := map[string]any{
    "prognosis_id": params.PrognosisID,
    "label":        params.Label,
    "answer":       params.Answer,
    "score":        params.Score,
  }

  builder := sq.Insert(questionTableName).
    SetMap(fields).
    Suffix(suffix)

  return getx[models.Question](ctx, s.conn, builder)
}

func (s *Storage) ChangeUserPassword(ctx context.Context, params models.ChangeUserPasswordParams) (*models.User, error) {
  suffix := returningSuffix(usersTableColumns)

  builder := sq.Update(usersTableName).
    Set("password", params.Password).
    Where(sq.Eq{"name": params.UserName}).
    Suffix(suffix)

  return getx[models.User](ctx, s.conn, builder)
}

func (s *Storage) GetUserByName(ctx context.Context, name string) (*models.User, error) {
  builder := sq.Select(usersTableColumns).
    From(usersTableName).
    Where(sq.Eq{"name": name})

  return getx[models.User](ctx, s.conn, builder)
}

func getx[T any](ctx context.Context, conn sqlscan.Querier, builder sq.Sqlizer) (*T, error) {
  query, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }

  dst := new(T)

  if err = sqlscan.Get(ctx, conn, dst, query, args...); err != nil {
    return nil, fmt.Errorf("sqlscan.Get: %w", err)
  }

  return dst, nil
}

func selectx[T any](ctx context.Context, conn sqlscan.Querier, builder sq.Sqlizer) ([]*T, error) {
  query, args, err := builder.ToSql()
  if err != nil {
    return nil, fmt.Errorf("builder.ToSql: %w", err)
  }

  var dst []*T

  if err = sqlscan.Select(ctx, conn, &dst, query, args...); err != nil {
    return nil, fmt.Errorf("sqlscan.Get: %w", err)
  }

  return dst, nil
}

func tableColumns(model any) string {
  typ := reflect.TypeOf(model)
  count := typ.NumField()

  columns := make([]string, 0, count)

  for idx := 0; idx < count; idx++ {
    if tag, ok := typ.Field(idx).Tag.Lookup("db"); ok {
      columns = append(columns, tag)
    }
  }

  return strings.Join(columns, ",")
}

func returningSuffix(columns string) string {
  return fmt.Sprintf("RETURNING %s", columns)
}
