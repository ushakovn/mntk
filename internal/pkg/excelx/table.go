package excelx

import (
  "bytes"
  "fmt"
  "math"
  "reflect"
  "strings"

  log "github.com/sirupsen/logrus"
  "github.com/xuri/excelize/v2"
)

const colWidthOffset = 0.3

var tableBorders = []excelize.Border{
  {Type: "top", Color: "#000000", Style: 1},
  {Type: "right", Color: "#000000", Style: 1},
  {Type: "bottom", Color: "#000000", Style: 1},
  {Type: "left", Color: "#000000", Style: 1},
}

type TableBuilder struct {
  tag   string
  title Title
  cells Cells
}

type Title struct {
  Color string
  Text  string
}

type Cells struct {
  Color  string
  Values any
}

func (t Title) defaults() Title {
  if t.Color == "" {
    t.Color = "#f9f9f9"
  }
  return t
}

func (b TableBuilder) Tag(value string) TableBuilder {
  b.tag = value
  return b
}

func (b TableBuilder) Title(value Title) TableBuilder {
  b.title = value.defaults()
  return b
}

func (b TableBuilder) Cells(values Cells) TableBuilder {
  b.cells = values
  return b
}

func (b TableBuilder) validate() error {
  if b.tag == "" {
    return fmt.Errorf("struct tag not specified")
  }
  return nil
}

func (b TableBuilder) Build() (string, error) {
  if err := b.validate(); err != nil {
    return "", fmt.Errorf("invalid builder: %w", err)
  }

  file := excelize.NewFile()

  defer func() {
    if err := file.Close(); err != nil {
      log.Errorf("excelx: file.Close: %v", err)
    }
  }()

  if err := file.SetSheetName("Sheet1", b.title.Text); err != nil {
    return "", fmt.Errorf("file.SetSheetName: %w", err)
  }

  sheetID, err := file.NewSheet(b.title.Text)
  if err != nil {
    return "", fmt.Errorf("file.NewSheet: %w", err)
  }
  file.SetActiveSheet(sheetID)

  if err = b.buildCells(file); err != nil {
    return "", fmt.Errorf("buildCells: %w", err)
  }

  buffer := new(bytes.Buffer)

  if err = file.Write(buffer); err != nil {
    return "", fmt.Errorf("file.SaveAs: %w", err)
  }

  return buffer.String(), nil
}

func (b TableBuilder) buildCells(file *excelize.File) error {
  sliceType := reflect.TypeOf(b.cells.Values)

  if sliceType.Kind() != reflect.Slice {
    return fmt.Errorf("not a slice argument")
  }
  elementType := sliceType.Elem()

  if elementType.Kind() == reflect.Ptr {
    elementType = elementType.Elem()
  }

  sliceValue := reflect.ValueOf(b.cells.Values)
  numElements := sliceValue.Len()

  if numElements != 0 {
    for i := 0; i < elementType.NumField(); i++ {
      fieldTag := elementType.Field(i).Tag.Get(b.tag)

      if fieldTag == `` || fieldTag == `-` {
        continue
      }

      cell := toCellIndex(i+1, 1)
      value := strings.TrimSpace(fieldTag)

      if err := file.SetCellValue(b.title.Text, cell, value); err != nil {
        return fmt.Errorf("file.SetCellValue: %w", err)
      }

      colWidth := float64(len(value)) + colWidthOffset
      col := toColIndex(i + 1)

      if err := file.SetColWidth(b.title.Text, col, col, colWidth); err != nil {
        return fmt.Errorf("file.SetColWidth: %w", err)
      }
    }

    styleID, err := file.NewStyle(&excelize.Style{
      Fill: excelize.Fill{
        Color:   []string{b.title.Color},
        Pattern: 1,
        Type:    "pattern",
      },
      Font: &excelize.Font{
        Family: "Arial",
        Bold:   true,
      },
      Border: tableBorders,
    })
    if err != nil {
      return fmt.Errorf("file.NewStyle: %w", err)
    }

    left := toCellIndex(1, 1)
    right := toCellIndex(elementType.NumField(), 1)

    if err = file.SetCellStyle(b.title.Text, left, right, styleID); err != nil {
      return fmt.Errorf("file.SetCellStyle: %w", err)
    }

    for j := 0; j < numElements; j++ {
      elementValue := sliceValue.Index(j)

      if elementValue.Kind() == reflect.Ptr {
        elementValue = elementValue.Elem()
      }

      for i := 0; i < elementType.NumField(); i++ {
        fieldTag := elementType.Field(i).Tag.Get(b.tag)

        if fieldTag == `` || fieldTag == `-` {
          continue
        }
        fieldValue := elementValue.Field(i)

        cell := toCellIndex(i+1, j+2)

        value := strings.TrimSpace(
          fmt.Sprintf("%v", fieldValue.Interface()),
        )

        if err = file.SetCellValue(b.title.Text, cell, value); err != nil {
          return fmt.Errorf("file.SetCellValue: %w", err)
        }
      }
    }

    styleID, err = file.NewStyle(&excelize.Style{
      Fill: excelize.Fill{
        Color:   []string{b.cells.Color},
        Pattern: 1,
        Type:    "pattern",
      },
      Font: &excelize.Font{
        Family: "Arial",
      },
      Border: tableBorders,
    })
    if err != nil {
      return fmt.Errorf("file.NewStyle: %w", err)
    }

    left = toCellIndex(1, 2)
    right = toCellIndex(elementType.NumField(), numElements+1)

    if err = file.SetCellStyle(b.title.Text, left, right, styleID); err != nil {
      return fmt.Errorf("file.SetCellStyle: %w", err)
    }
  }

  return nil
}

func toCellIndex(col, row int) string {
  out := toColIndex(col)

  return fmt.Sprint(out, row)
}

func toColIndex(col int) string {
  buffer := make([]rune, int(math.Floor(
    math.Log(float64(25*(col+1)))/math.Log(26)),
  ))

  for i := len(buffer) - 1; i >= 0; i-- {
    col--
    buffer[i] = rune('A' + col%26)
    col /= 26
  }
  out := string(buffer)

  return out
}
