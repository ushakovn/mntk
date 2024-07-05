package authx

import (
  "context"
  "errors"
  "fmt"
  "mntk/internal/models"
  "net/http"
)

const userName = "admin"

var (
  ErrWrongOldPassword = errors.New("wrong password")
  ErrNewPasswordSame  = errors.New("new password same")
)

type Config struct {
  Storage Storage
}

type Handler struct {
  storage Storage
}

type Storage interface {
  ChangeUserPassword(ctx context.Context, params models.ChangeUserPasswordParams) (*models.User, error)
  GetUserByName(ctx context.Context, name string) (*models.User, error)
}

func NewHandler(config Config) *Handler {
  return &Handler{
    storage: config.Storage,
  }
}

func (h *Handler) ChangePassword(ctx context.Context, oldPassword, newPassword string) error {
  admin, err := h.storage.GetUserByName(ctx, userName)
  if err != nil {
    return fmt.Errorf("storage.GetUserByName: %w", err)
  }

  if admin.Password != oldPassword {
    return ErrWrongOldPassword
  }

  if admin.Password == newPassword {
    return ErrNewPasswordSame
  }

  _, err = h.storage.ChangeUserPassword(ctx, models.ChangeUserPasswordParams{
    UserName: userName,
    Password: newPassword,
  })
  if err != nil {
    return fmt.Errorf("storage.ChangeUserPassword: %w", err)
  }

  return nil
}

func (h *Handler) Middleware(handler http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    user, password, ok := r.BasicAuth()
    if !ok {
      writeBasicAuthRetryResp(w, "Не введены учетные данные")
      return
    }

    admin, err := h.storage.GetUserByName(ctx, userName)
    if err != nil {
      http.Error(w, "Непредвиденная ошибка", http.StatusInternalServerError)
      return
    }

    if user != admin.Name || password != admin.Password {
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
