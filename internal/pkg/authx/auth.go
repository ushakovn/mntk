package authx

import (
  "context"
  "errors"
  "fmt"
  "mntk/internal/models/user"
  "net/http"
)

const userName = "admin"

var (
  ErrWrongPassword   = errors.New("wrong password")
  ErrNewPasswordSame = errors.New("new password same")
)

type Dependencies struct {
  User UserProvider
}

type Provider struct {
  user UserProvider
}

type UserProvider interface {
  GetByName(ctx context.Context, name string) (*user.User, error)
  ChangePassword(ctx context.Context, params user.ChangePasswordParams) (*user.User, error)
}

func NewProvider(deps Dependencies) *Provider {
  return &Provider{
    user: deps.User,
  }
}

func (h *Provider) ChangePassword(ctx context.Context, oldPassword, newPassword string) error {
  admin, err := h.user.GetByName(ctx, userName)
  if err != nil {
    return fmt.Errorf("user.GetByName: %w", err)
  }

  if admin.Password != oldPassword {
    return ErrWrongPassword
  }

  if admin.Password == newPassword {
    return ErrNewPasswordSame
  }

  _, err = h.user.ChangePassword(ctx, user.ChangePasswordParams{
    Name:     userName,
    Password: newPassword,
  })
  if err != nil {
    return fmt.Errorf("user.ChangePassword: %w", err)
  }

  return nil
}

func (h *Provider) Middleware(handler http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    name, password, ok := r.BasicAuth()
    if !ok {
      writeBasicAuthRetryResp(w, "Не введены учетные данные")
      return
    }

    admin, err := h.user.GetByName(ctx, userName)
    if err != nil {
      http.Error(w, "Непредвиденная ошибка", http.StatusInternalServerError)
      return
    }

    if name != admin.Name || password != admin.Password {
      writeBasicAuthRetryResp(w, "Некорректные учетные данные")
      return
    }

    handler(w, r)
  }
}

func writeBasicAuthRetryResp(w http.ResponseWriter, msg string) {
  w.Header().Add(
    "WWW-Authenticate",
    "Basic realm=admin",
  )
  http.Error(w, msg, http.StatusUnauthorized)
}
