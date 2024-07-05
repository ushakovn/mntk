package pagination

type Pagination struct {
  Limit  uint64
  Offset uint64
}

func (p Pagination) IsValid() bool {
  return p.Limit != 0 || p.Offset != 0
}
