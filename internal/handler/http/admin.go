package http

import (
  "errors"
  "fmt"
  "mntk/internal/pkg/authx"
  "mntk/internal/pkg/httpx"
  "mntk/internal/templates/forms"
  "net/http"
)

type ChangeAdminPasswordRequest struct {
  Current string `json:"current"`
  New     string `json:"new"`
}

func (h *Handler) ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[ChangeAdminPasswordRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }

  ctx := r.Context()

  if err = h.providers.Auth.ChangePassword(ctx, req.Current, req.New); err != nil {
    switch {
    case errors.Is(err, authx.ErrWrongPassword):
      http.Error(w, "Введен неверный старый пароль", http.StatusBadRequest)

    case errors.Is(err, authx.ErrNewPasswordSame):
      http.Error(w, "Новый пароль должен быть отличен от старого", http.StatusBadRequest)

    default:
      httpx.WriteInternalError(w, r, fmt.Errorf("auth.ChangePassword: %w", err))
    }
    return
  }

  w.WriteHeader(http.StatusOK)
}

func (h *Handler) ChangeAdminPasswordForm(w http.ResponseWriter, r *http.Request) {
  httpx.WriteBytes(w, r, forms.ChangeAdminPassword())
}

func (h *Handler) AdminLogout(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "", http.StatusUnauthorized)
}

func (h *Handler) AdminForm(w http.ResponseWriter, r *http.Request) {
  httpx.WriteBytes(w, r, forms.Admin())
}
