package httpx

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"

  log "github.com/sirupsen/logrus"
)

func ReadRequest[T any](r *http.Request) (*T, error) {
  var payload []byte

  switch r.Method {

  case http.MethodPost:
    defer func() {
      _ = r.Body.Close()
    }()

    p, err := io.ReadAll(r.Body)
    if err != nil {
      return nil, fmt.Errorf("cannot read request payload: %v", err)
    }
    payload = p

  case http.MethodOptions:
    return nil, nil

  default:
    return nil, fmt.Errorf("unsupported method")
  }

  req := new(T)

  if err := json.Unmarshal(payload, req); err != nil {
    return nil, fmt.Errorf("cannot parse request payload to req")
  }

  return req, nil
}

func WriteBytes(w http.ResponseWriter, r *http.Request, p []byte) {
  if _, err := w.Write(p); err != nil {
    WriteInternalError(w, r, fmt.Errorf("writer.Write: %w", err))
  }
}

func WriteResponse(w http.ResponseWriter, r *http.Request, p any) {
  b, err := json.MarshalIndent(p, "", "  ")
  if err != nil {
    WriteInternalError(w, r, fmt.Errorf("cannot marshal indent response: %w", err))
    return
  }

  w.Header().Add("Content-Type", "application/json")

  if _, err = w.Write(b); err != nil {
    WriteInternalError(w, r, fmt.Errorf("cannot write to message writer: %v", err))
    return
  }
}

func WriteRequestError(w http.ResponseWriter, r *http.Request, err error) {
  log.Infof("error in request %s: %v", r.RequestURI, err)
  http.Error(w, err.Error(), http.StatusBadRequest)
}

func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
  log.Infof("internal error in request %s: %v", r.RequestURI, err)
  http.Error(w, err.Error(), http.StatusInternalServerError)
}

func HandleHealthRequest(w http.ResponseWriter, _ *http.Request) {
  w.WriteHeader(http.StatusOK)

  if _, err := w.Write([]byte("/ok")); err != nil {
    log.Printf("cannot write to response writer: %v", err)
  }
}
