package prognosis

import (
  "mntk/internal/models/pagination"
  "mntk/internal/models/sort"
  "mntk/internal/pkg/sqlx"
  "time"
)

const (
  CreatedAtPrognosisSortField SortField = "created_at"
  ScorePrognosisSortField     SortField = "score"
)

type SortField string

type ID int64

type Prognosis struct {
  ID     ID      `db:"id"`
  Score  float64 `db:"score"`
  Result string  `db:"result"`

  PatientName  string    `db:"patient_name"`
  PatientBirth time.Time `db:"patient_birth"`

  CreatedAt time.Time `db:"created_at"`
}

type CreateParams struct {
  Score  float64
  Result string

  PatientName  string
  PatientBirth time.Time
}

type ListParams struct {
  Pagination pagination.Pagination
  Sort       Sort
}

type Sort struct {
  Field SortField
  Order sort.Order
}

func (p Sort) IsValid() bool {
  return sqlx.IsValidOrderBy(p.Field, p.Order)
}

func (p Sort) OrderBy() string {
  return sqlx.OrderBy(p.Field, p.Order)
}
