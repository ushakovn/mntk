package models

import (
  "fmt"
  "time"
)

const (
  CreatedAtPrognosisSortField PrognosisSortField = "created_at"
  ScorePrognosisSortField     PrognosisSortField = "score"
)

const (
  AscSortOrder  SortOrder = "asc"
  DescSortOrder SortOrder = "desc"
)

type PrognosisID int64

type Prognosis struct {
  ID          PrognosisID `db:"id"`
  Score       float64     `db:"score"`
  Result      string      `db:"result"`
  PatientName string      `db:"patient_name"`
  CreatedAt   time.Time   `db:"created_at"`
}

type QuestionID int64

type Question struct {
  ID          QuestionID  `db:"id"`
  PrognosisID PrognosisID `db:"prognosis_id"`
  Label       string      `db:"label"`
  Answer      string      `db:"answer"`
  Score       float64     `db:"score"`
  CreatedAt   time.Time   `db:"created_at"`
}

type CreatePrognosisParams struct {
  Score       float64
  Result      string
  PatientName *string
}

type ListPrognosisParams struct {
  Pagination Pagination
  Sort       PrognosisSort
}

type PrognosisSortField string

type PrognosisSort struct {
  Field PrognosisSortField
  Order SortOrder
}

type CreateQuestionParams struct {
  PrognosisID PrognosisID
  Label       string
  Answer      string
  Score       float64
}

type User struct {
  ID       int64  `db:"id"`
  Name     string `db:"name"`
  Password string `db:"password"`
}

type ChangeUserPasswordParams struct {
  UserName string
  Password string
}

type SortOrder string

type Pagination struct {
  Limit  uint64
  Offset uint64
}

func (p Pagination) IsValid() bool {
  return p.Limit != 0 || p.Offset != 0
}

func (p PrognosisSort) IsValid() bool {
  return p.Field != "" || p.Order != ""
}

func (p PrognosisSort) OrderBy() string {
  return fmt.Sprintf("%s %s",
    p.Field,
    p.Order,
  )
}
