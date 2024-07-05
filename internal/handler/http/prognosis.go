package http

import (
  "fmt"
  "mntk/internal/models/pagination"
  "mntk/internal/models/prognosis"
  "mntk/internal/models/sort"
  "mntk/internal/pkg/excelx"
  "mntk/internal/pkg/htmlx"
  "mntk/internal/pkg/httpx"
  "mntk/internal/templates/forms"
  "net/http"
  "time"
)

type Prognosis struct {
  ID int64 `json:"id" output:"Номер записи"`

  Score  float64 `json:"score" output:"Набранное количество баллов"`
  Result string  `json:"result" output:"Полученный прогноз"`

  PatientName  string `json:"patient_name" output:"Имя пациента"`
  PatientBirth string `json:"patient_birth" output:"Дата рождения"`

  CreatedAt string `json:"created_at" output:"Время прохождения опроса"`
}

type CreatePrognosisRequest struct {
  Score  float64 `json:"score"`
  Result string  `json:"result"`

  Patient struct {
    Name  string `json:"name"`
    Birth string `json:"birth"`
  }
}

type CreatePrognosisResponse struct {
  Prognosis Prognosis `json:"prognosis"`
}

func (h *Handler) CreatePrognosis(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[CreatePrognosisRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }
  ctx := r.Context()

  p, err := h.providers.Prognosis.Create(ctx, prognosis.CreateParams{
    Score:  req.Score,
    Result: req.Result,

    PatientName: req.Patient.Name,

    PatientBirth: func() time.Time {
      date, _ := time.Parse(reqDateFormat, req.Patient.Birth)

      return date
    }(),
  })
  if err != nil {
    httpx.WriteInternalError(w, r, fmt.Errorf("providers.Prognosis.Create: %w", err))
    return
  }

  resp := Prognosis{
    ID: int64(p.ID),

    Score:  p.Score,
    Result: p.Result,

    PatientName:  p.PatientName,
    PatientBirth: p.PatientBirth.Format(respDateFormat),

    CreatedAt: p.CreatedAt.Format(respTimeFormat),
  }

  httpx.WriteResponse(w, r, CreatePrognosisResponse{
    Prognosis: resp,
  })
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

func (h *Handler) ListPrognosis(w http.ResponseWriter, r *http.Request) {
  req, err := httpx.ReadRequest[ListPrognosisRequest](r)
  if err != nil {
    httpx.WriteRequestError(w, r, fmt.Errorf("invalid request: %w", err))
    return
  }
  ctx := r.Context()

  ps, err := h.providers.Prognosis.List(ctx, prognosis.ListParams{
    Pagination: pagination.Pagination{
      Limit:  req.Limit,
      Offset: req.Offset,
    },
    Sort: prognosis.Sort{
      Field: prognosis.SortField(req.SortField),
      Order: sort.Order(req.SortOrder),
    },
  },
  )
  if err != nil {
    httpx.WriteInternalError(w, r, fmt.Errorf("providers.Prognosis.List: %w", err))
    return
  }

  resp := make([]*Prognosis, 0, len(ps))

  for _, p := range ps {
    resp = append(resp, &Prognosis{
      ID: int64(p.ID),

      Score:  p.Score,
      Result: p.Result,

      PatientName:  p.PatientName,
      PatientBirth: p.PatientBirth.Format(respDateFormat),

      CreatedAt: p.CreatedAt.Format(respTimeFormat),
    })
  }

  switch r.Header.Get("Content-Type") {

  case "text/html":
    resp, err := buildListPrognosisHtml(resp)
    if err != nil {
      httpx.WriteInternalError(w, r, fmt.Errorf("buildListPrognosisHtml: %w", err))
    }
    httpx.WriteBytes(w, r, []byte(resp))

  case "application/xls":
    resp, err := buildListPrognosisExcel(resp)
    if err != nil {
      httpx.WriteInternalError(w, r, fmt.Errorf("buildListPrognosisExcel: %w", err))
    }

    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", `attachment;filename="file.xlsx"`)
    w.Header().Set("Content-Transfer-Encoding", "binary")

    httpx.WriteBytes(w, r, []byte(resp))

  case "application/json":
    httpx.WriteResponse(w, r, ListPrognosisResponse{
      Prognosis: resp,
    })

  default:
    httpx.WriteRequestError(w, r, fmt.Errorf("unsupported Content-Type"))
  }
}

func buildListPrognosisExcel(cells []*Prognosis) (string, error) {
  var builder excelx.TableBuilder

  builder = builder.
    Tag("output").
    Title(excelx.Title{
      Color: "#e9eff6",
      Text:  "📑 Результаты опросов",
    }).
    Cells(excelx.Cells{
      Color:  "#ffffff",
      Values: cells,
    })

  str, err := builder.Build()
  if err != nil {
    return "", fmt.Errorf("builder.Build: %w", err)
  }

  return str, nil
}

func buildListPrognosisHtml(cells []*Prognosis) (string, error) {
  var builder htmlx.TableBuilder

  builder = builder.
    Tag("output").
    Title("Результаты опросов").
    Header("📑 Результаты опросов").
    Cells(cells).
    RedirectButton("Вернуться назад", "/admin",
      htmlx.TopButtonPosition,
    ).
    ScriptButton("Выгрузить в Excel", `
        try {
            let resp = await fetch("/prognosis/list", {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/xls'
                },
                body: JSON.stringify(
                    {
                        "sort_order": "desc",
                        "sort_field": "created_at"
                    })
            })

            let blob = await resp.blob();
            let url = window.URL.createObjectURL(blob);

            let a = document.createElement('a');

            a.href = url;
            a.download = "export.xlsx";

            document.body.appendChild(a);

            a.click();
            a.remove();  

        } catch (error) {
            console.error(error.message);
        }
      `,
      htmlx.TopButtonPosition)

  str, err := builder.Build()
  if err != nil {
    return "", fmt.Errorf("builder.Build: %w", err)
  }

  return str, nil
}

func (h *Handler) PrognosisForm(w http.ResponseWriter, r *http.Request) {
  httpx.WriteBytes(w, r, forms.Prognosis())
}
