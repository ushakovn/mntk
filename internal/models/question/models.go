package question

import (
  "mntk/internal/models/prognosis"
  "time"
)

type ID int64

type Question struct {
  ID          ID           `db:"id"`
  PrognosisID prognosis.ID `db:"prognosis_id"`
  Label       string       `db:"label"`
  Answer      string       `db:"answer"`
  Score       float64      `db:"score"`
  CreatedAt   time.Time    `db:"created_at"`
}

type CreateParams struct {
  PrognosisID prognosis.ID
  Label       string
  Answer      string
  Score       float64
}
