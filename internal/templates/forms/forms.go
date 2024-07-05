package forms

import "text/template"

var funcMap = template.FuncMap{
  "Inc": func(a int) int {
    return a + 1
  },
}

func init() {
  buf, err := buildPrognosis()
  if err != nil {
    panic(err)
  }
  prognosisBuf = buf
}

func Prognosis() []byte {
  return prognosisBuf
}

func Admin() []byte {
  return adminBuf
}

func ChangeAdminPassword() []byte {
  return changeAdminPasswordBuf
}
