package http

import (
  "context"
  "errors"
  "fmt"
  "mntk/internal/models"
  "mntk/internal/pkg/authx"
  "mntk/internal/pkg/htmlx"
  "mntk/internal/pkg/httpx"
  "mntk/internal/templates/forms"
  "net/http"

  log "github.com/sirupsen/logrus"
)

const respTimeFormat = "02-01-2006 15:04:05"

type Config struct {
  Port    int
  Storage Storage
  Auth    AuthHandler
}

type Handler struct {
  port    int
  storage Storage
  auth    AuthHandler
}

type Storage interface {
  ListPrognosis(ctx context.Context, params models.ListPrognosisParams) ([]*models.Prognosis, error)
  CreatePrognosis(ctx context.Context, params models.CreatePrognosisParams) (*models.Prognosis, error)
  CreateQuestion(ctx context.Context, params models.CreateQuestionParams) (*models.Question, error)
}

type AuthHandler interface {
  Middleware(handler http.HandlerFunc) http.HandlerFunc
  ChangePassword(ctx context.Context, oldPassword, newPassword string) error
}

type CreatePrognosisRequest struct {
  Score       float64 `json:"score"`
  Result      string  `json:"result"`
  PatientName *string `json:"patient_name,omitempty"`
}

type CreatePrognosisResponse struct {
  Prognosis Prognosis `json:"prognosis"`
}

type ListPrognosisRequest struct {
  Limit     uint64 `json:"limit"`
  Offset    uint64 `json:"offset"`
  SortField string `json:"sort_field"`
  SortOrder string `json:"sort_order"`
}

type ListPrognosisResponse struct {
  Prognosis []*Prognosis `json:"prognosis"`
}

type Prognosis struct {
  ID          int64   `json:"id" html:"Номер записи"`
  Score       float64 `json:"score" html:"Набранное количество баллов"`
  Result      string  `json:"result" html:"Полученный прогноз"`
  PatientName string  `json:"patient_name" html:"-"`
  CreatedAt   string  `json:"created_at" html:"Время прохождения опроса"`
}

type ChangeAdminPasswordRequest struct {
  Current string `json:"current"`
  New     string `json:"new"`
}

func NewHandler(config Config) *Handler {
  return &Handler{
    port:    config.Port,
    storage: config.Storage,
    auth:    config.Auth,
  }
}

func (h *Handler) ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[ChangeAdminPasswordRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }

  ctx := r.Context()

  if err = h.auth.ChangePassword(ctx, req.Current, req.New); err != nil {
    switch {
    case errors.Is(err, authx.ErrWrongOldPassword):
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

func (h *Handler) PrognosisForm(w http.ResponseWriter, r *http.Request) {
  httpx.WriteBytes(w, r, forms.Prognosis())
}

func (h *Handler) IndexRedirect(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/prognosis", http.StatusPermanentRedirect)
}

func (h *Handler) CreatePrognosis(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[CreatePrognosisRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }
  ctx := r.Context()

  p, err := h.storage.CreatePrognosis(ctx, models.CreatePrognosisParams{
    Score:       req.Score,
    Result:      req.Result,
    PatientName: req.PatientName,
  })
  if err != nil {
    httpx.WriteInternalError(w, r, fmt.Errorf("storage.CreatePrognosis: %w", err))
    return
  }

  respP := Prognosis{
    ID:          int64(p.ID),
    Score:       p.Score,
    Result:      p.Result,
    PatientName: p.PatientName,
    CreatedAt:   p.CreatedAt.Format(respTimeFormat),
  }

  httpx.WriteResponse(w, r, CreatePrognosisResponse{
    Prognosis: respP,
  })
}

func (h *Handler) ListPrognosis(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[ListPrognosisRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }
  ctx := r.Context()

  ps, err := h.storage.ListPrognosis(ctx,
    models.ListPrognosisParams{
      Pagination: models.Pagination{
        Limit:  req.Limit,
        Offset: req.Offset,
      },
      Sort: models.PrognosisSort{
        Field: models.PrognosisSortField(req.SortField),
        Order: models.SortOrder(req.SortOrder),
      },
    },
  )
  if err != nil {
    httpx.WriteInternalError(w, r, fmt.Errorf("storage.ListPrognosis: %w", err))
    return
  }

  respP := make([]*Prognosis, 0, len(ps))

  for _, p := range ps {
    respP = append(respP, &Prognosis{
      ID:          int64(p.ID),
      Score:       p.Score,
      Result:      p.Result,
      PatientName: p.PatientName,
      CreatedAt:   p.CreatedAt.Format(respTimeFormat),
    })
  }

  if typ := r.Header.Get("Content-Type"); typ == "text/html" {
    resp, err2 := htmlx.MakeTable(
      " Результаты опросов",
      "📑 Результаты опросов",
      "/admin",
      respP,
    )
    if err2 != nil {
      httpx.WriteInternalError(w, r, fmt.Errorf("htmlx.MakeTable: %w", err))
    }
    httpx.WriteBytes(w, r, resp)
    return
  }

  httpx.WriteResponse(w, r, ListPrognosisResponse{
    Prognosis: respP,
  })
}

func (h *Handler) Run() {
  h.registerHandles()

  log.Infof("run http handler on: http://localhost:%d", h.port)
  addr := fmt.Sprintf(":%d", h.port)

  if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatalf("server listening error: %v", err)
  }
}

func (h *Handler) registerHandles() {
  registerHandle("/prognosis/create", h.CreatePrognosis)
  registerHandle("/prognosis/list", h.ListPrognosis)
  registerHandle("/prognosis", h.PrognosisForm)

  registerHandle("/admin", h.auth.Middleware(h.AdminForm))
  registerHandle("/admin/password", h.auth.Middleware(h.ChangeAdminPasswordForm))
  registerHandle("/admin/password/change", h.ChangeAdminPassword)
  registerHandle("/admin/logout", h.AdminLogout)

  registerHandle("/health", httpx.HandleHealthRequest)

  registerHandle("/", h.IndexRedirect)
  registerAssetsHandle()
}

func registerHandle(route string, handler func(w http.ResponseWriter, r *http.Request)) {
  http.Handle(route, http.HandlerFunc(handler))
}

func registerAssetsHandle() {
  const (
    dir   = "internal/templates/forms/assets"
    route = "/assets/"
  )
  fs := http.FileServer(http.Dir(dir))
  http.Handle(route, http.StripPrefix(route, fs))
}
