package htmlx

import (
  "fmt"
  "reflect"
  "strings"
)

const (
  TopButtonPosition    TableButtonPosition = 0
  BottomButtonPosition TableButtonPosition = 1
)

type TableButtonPosition int

type TableBuilder struct {
  tag     string
  title   string
  header  string
  cells   any
  buttons []tableButton
}

type tableButton struct {
  Text     string
  Script   string
  Position TableButtonPosition
}

func (b TableBuilder) Tag(value string) TableBuilder {
  b.tag = value
  return b
}

func (b TableBuilder) Title(value string) TableBuilder {
  b.title = value
  return b
}

func (b TableBuilder) Header(value string) TableBuilder {
  b.header = value
  return b
}

func (b TableBuilder) Cells(values any) TableBuilder {
  b.cells = values
  return b
}

func (b TableBuilder) RedirectButton(text, path string, pos TableButtonPosition) TableBuilder {
  index := len(b.buttons)

  script := fmt.Sprintf(`
    <script>
      document.getElementById('button-%d').addEventListener('click', function() {
        window.location.href = '%s';
      });
    </script>
  `, index, path)

  b.buttons = append(b.buttons, tableButton{
    Text:     text,
    Script:   script,
    Position: pos,
  })

  return b
}

func (b TableBuilder) ScriptButton(text, script string, pos TableButtonPosition) TableBuilder {
  index := len(b.buttons)

  script = fmt.Sprintf(`
    <script>
      document.getElementById('button-%d').addEventListener('click', async function() {
        %s
      });
    </script>
  `, index, script)

  b.buttons = append(b.buttons, tableButton{
    Text:     text,
    Script:   script,
    Position: pos,
  })

  return b
}

func (b TableBuilder) validate() error {
  if b.tag == "" {
    return fmt.Errorf("struct tag not specified")
  }
  if b.header == "" {
    return fmt.Errorf("table header not specified")
  }
  if b.title == "" {
    return fmt.Errorf("table title not specified")
  }
  return nil
}

func (b TableBuilder) Build() (string, error) {
  if err := b.validate(); err != nil {
    return "", fmt.Errorf("invalid builder: %w", err)
  }

  var sb strings.Builder

  str := fmt.Sprintf(`
    <body>

      <div class="container">
      <img class="ocular" src="assets/ocular.svg">

      <h5>%s</h5>
      <hr/>
`, b.header)

  sb.WriteString(str)

  for index, button := range b.buttons {
    if button.Position == TopButtonPosition {
      str := buildButton(index, button.Text, button.Script)
      sb.WriteString(str)
    }
  }

  sb.WriteString(`
      </div>
      <div class="container table-container">
`)

  str, err := buildCells(b.tag, b.cells)
  if err != nil {
    return "", fmt.Errorf("buildCells: %w", err)
  }
  sb.WriteString(str)

  for index, button := range b.buttons {
    if button.Position == BottomButtonPosition {
      str = buildButton(index, button.Text, button.Script)
      sb.WriteString(str)
    }
  }

  sb.WriteString(`
      </div>

    </body>
`)

  str = fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="ru">
    
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
    
        <link 
          href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" 
          rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" 
          crossorigin="anonymous"
        >
    
        <title>%s</title>
    
        <style>
            body {
                font-family: Arial, sans-serif;
                background-image: url('/assets/background.svg');
                background-size: cover;
                background-repeat: no-repeat;
                /*display: flex;*/
                justify-content: center;
                align-items: center;
                background-color: #bfcfde;
                /* Центрирование */
                margin: 0;
                /* height: 100vh; */
                display: grid;
                place-items: center;

                margin-left: auto;
                margin-right: auto;
            }
    
            .container {
                width: 100%%;
                max-width: 800px;
                /* Без центрирования */
                margin: 10px;
                margin-top: 20px;
                padding: 20px;
                background-color: rgba(255, 255, 255, 1);
                border-radius: 10px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
    
            h5 {
                text-align: center;
                margin-top: 4px;
            }
    
            button {
                width: 100%%;
                padding: 10px;
                margin-top: 10px;
                background-color: #95abc3;
                color: #fff;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                font-size: 18px;
            }
    
            button:hover {
                background-color: #7d9bb7;
            }

            table {
                font-size: 12px;
            }

            .table-container {
                width: auto;            
            }

            .ocular {
                display: none;
            }

            /* Стили для мобильных и планшет устройств */
            @media (max-width: 815px) {
                .wrap {
                    border-radius: 3px;
                    margin-top: 18px;
                    margin-bottom: 20px;
                    width: auto;
                    max-width: 340px;
                    padding: 5px;
                    background-color: rgba(255, 255, 255, 1);
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    font-size: 13px;
                }
                .container {
                    display: block;
                    padding: 10px;
                    font-size: 14px;
                    width: auto;
                    margin: inherit;
                }
                .table-container {
                    width: auto;            
                }
                .container button {
                    display: block;
                    padding: 6px;
                    border-radius: 3px;
                    font-size: 13px;
                    width: 94%%;
                    margin: 8px auto;
                }
                hr {
                    width: 92%%;
                    display: block;
                    margin: 14px auto;
                }
    
                .question {
                    border: 1px solid #ccc;
                    padding: 10px 5px 5px 10px;
                }
                .input-patient {
                    font-size: 12px;
                    margin-bottom: 8px;
                    min-width: 280px;
                }
                label {
                    margin-bottom: 6px;
                }
                .question label {
                    font-size: 12px;
                }
                .options label {
                    font-size: 12px;
                }
                h3 {
                    text-align: center;
                    margin: 8px auto 8px;
                    font-size: 13px;
                    font-weight: 300;
                    display: block;
                    max-width: 80%%;
                }
                h5 {
                    font-size: 14px;
                }
                .floating-button {
                    display: none;
                }
                .error-message {
                    margin-top: 2px;
                    margin-bottom: 2px;
                    font-size: 12px;
                }
                body {
                    height: auto;
                }
    
                .result {
                    display: block;
                    width: 93%%;
                    border: 2px solid #ccc;
                    margin-left: auto;
                    margin-right: auto;
                    margin-top: 12px;
                }
                .result-table {
                    display: none;
                }
                .result-table-mobile {
                    display: table;
                    font-size: inherit;
                    border: 1px solid #dddddd;
                    width: 320px;
                    margin-left: auto;
                    margin-right: auto;
                }
                .result-p {
                    font-size: 12px;
                }
    
                /* Кнопка администратора */
                .floating-button {
                    display: block;
                    font-size: 14px;
                    padding: 6px;
                    max-width: 35px;
                    right: 30px;
                    bottom: 80px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                }
    
                body {
                    display: block;
                    margin-left: auto;
                    margin-right: auto;
                    margin-top: 20px;
                    background-size: inherit;
                }

                table {
                    display: table;
                    font-size: 6px;
                    border: 1px solid #dddddd;
                    margin-left: auto;
                    margin-right: auto;
                }
            }

            /* Стили для десктоп */
            @media (max-width: 1920px) {
                .wrap {
                    font-size: 14px;
                    width: inherit;
                    /*margin-top: 10px;*/
                }
                label {
                    margin-bottom: 8px;
                    font-size: 14px;
                }
                input {
                    font-size: 14px;
                }
                .input-patient {
                    font-size: 14px;
                }
                button {
                    font-size: 14px;
                }
                .question {
                    margin-bottom: 8px;
                    padding: 8px 8px 8px 14px;
                }
                .section {
                    margin: 8px 6px;
                }
                h3 {
                    margin-bottom: 12px;
                }
                .options {
                    font-size: 13px;
                }
                .error-message {
                    font-size: 13px;
                    margin-top: 0;
    
                }
                .input-patient {
                    margin-bottom: 8px;
                    max-width: 270px;
                }
                h5 {
                    font-size: 16px;
                }
                hr {
                    margin: 10px;
                }
            }
        </style>
    </head>

      %s

    </html>
    `,
    b.title,
    sb.String(),
  )

  return str, nil
}

func buildCells(tag string, cells any) (string, error) {
  sliceType := reflect.TypeOf(cells)

  if sliceType.Kind() != reflect.Slice {
    return "", fmt.Errorf("not a slice argument")
  }
  elementType := sliceType.Elem()

  if elementType.Kind() == reflect.Ptr {
    elementType = elementType.Elem()
  }

  sliceValue := reflect.ValueOf(cells)
  numElements := sliceValue.Len()

  var (
    headers strings.Builder
    rows    strings.Builder
  )

  if numElements != 0 {
    headers.WriteString(`<tr>`)

    for i := 0; i < elementType.NumField(); i++ {

      fieldTag := elementType.Field(i).Tag.Get(tag)

      if fieldTag == `` || fieldTag == `-` {
        continue
      }
      headers.WriteString(fmt.Sprintf(`<th>%s</th>`, fieldTag))
    }
    headers.WriteString(`</tr>`)

    for j := 0; j < numElements; j++ {
      rows.WriteString(`<tr>`)
      elementValue := sliceValue.Index(j)

      if elementValue.Kind() == reflect.Ptr {
        elementValue = elementValue.Elem()
      }

      for i := 0; i < elementType.NumField(); i++ {

        fieldTag := elementType.Field(i).Tag.Get(tag)

        if fieldTag == `` || fieldTag == `-` {
          continue
        }
        fieldValue := elementValue.Field(i)

        rows.WriteString(fmt.Sprintf(`<td>%v</td>`, fieldValue.Interface()))
      }
      rows.WriteString(`</tr>`)
    }
  }

  str := fmt.Sprintf(`
    <table class="table table-bordered">
      %s
      %s
    </table>
  `, headers.String(), rows.String())

  return str, nil
}

func buildButton(index int, text, script string) string {
  str := fmt.Sprintf(`
        <button type="button" id="button-%d">
          %s
        </button>
        %s
      `, index, text, script)

  return str
}
