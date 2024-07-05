package http

import (
  "fmt"
  "mntk/internal/pkg/httpx"
  "net/http"

  log "github.com/sirupsen/logrus"
)

const (
  reqDateFormat  = "2006-01-02"
  respDateFormat = "02.01.2006"
  respTimeFormat = "02-01-2006 15:04:05"
)

type Handler struct {
  config    Config
  providers Providers
}

type Config struct {
  Port int
}

func NewHandler(config Config, deps Dependencies) *Handler {
  return &Handler{
    config:    config,
    providers: deps.Providers,
  }
}

func (h *Handler) Run() {
  h.registerHandles()

  log.Infof("run http handler on: http://localhost:%d", h.config.Port)
  addr := fmt.Sprintf(":%d", h.config.Port)

  if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatalf("server listening error: %v", err)
  }
}

func (h *Handler) IndexRedirect(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/prognosis", http.StatusPermanentRedirect)
}

func (h *Handler) registerHandles() {
  registerHandle("/prognosis/create", h.CreatePrognosis)
  registerHandle("/prognosis/list", h.ListPrognosis)

  registerHandle("/prognosis", h.PrognosisForm)

  registerHandle("/admin", h.providers.Auth.Middleware(h.AdminForm))
  registerHandle("/admin/password", h.providers.Auth.Middleware(h.ChangeAdminPasswordForm))

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
