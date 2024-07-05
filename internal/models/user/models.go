package user

type User struct {
  ID       int64  `db:"id"`
  Name     string `db:"name"`
  Password string `db:"password"`
}

type ChangePasswordParams struct {
  Name     string
  Password string
}
