package sqlx

import (
  "context"
  "fmt"
  "reflect"
  "strings"

  sq "github.com/Masterminds/squirrel"
  "github.com/georgysavva/scany/v2/sqlscan"
)

func Get[T any](ctx context.Context, conn sqlscan.Querier, builder sq.Sqlizer) (*T, error) {
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

func Select[T any](ctx context.Context, conn sqlscan.Querier, builder sq.Sqlizer) ([]*T, error) {
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

func Columns(model any) string {
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

func ReturningSuffix(columns string) string {
  return fmt.Sprintf("RETURNING %s", columns)
}

func OrderBy[F, O ~string](field F, order O) string {
  return fmt.Sprintf("%s %s", field, order)
}

func IsValidOrderBy[F, O ~string](field F, order O) bool {
  return field != "" && order != ""
}
