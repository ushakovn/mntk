package htmlx

import (
  "fmt"
  "reflect"
  "strings"
)

const TagKey = `html`

func MakeTable(title, description, redirectBack string, slice any) ([]byte, error) {
  sliceType := reflect.TypeOf(slice)

  if sliceType.Kind() != reflect.Slice {
    return nil, fmt.Errorf("not a slice argument")
  }
  elementType := sliceType.Elem()

  if elementType.Kind() == reflect.Ptr {
    elementType = elementType.Elem()
  }

  sliceValue := reflect.ValueOf(slice)
  numElements := sliceValue.Len()

  var (
    headers strings.Builder
    rows    strings.Builder
  )

  if numElements != 0 {
    headers.WriteString(`<tr>`)

    for i := 0; i < elementType.NumField(); i++ {

      fieldTag := elementType.Field(i).Tag.Get(TagKey)

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

        fieldTag := elementType.Field(i).Tag.Get(TagKey)

        if fieldTag == `` || fieldTag == `-` {
          continue
        }
        fieldValue := elementValue.Field(i)

        rows.WriteString(fmt.Sprintf(`<td>%v</td>`, fieldValue.Interface()))
      }
      rows.WriteString(`</tr>`)
    }
  }

  html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ru">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">

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
            /*margin: 0;*/
            /*height: 100vh;*/
            display: grid;
            place-items: center;
        }

        .container {
            width: 100%%;
            max-width: 720px;
            /* Без центрирования */
            margin: 30px;
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
    </style>
</head>

<body>

  <div class="container">
    <button type="button" id="backBtn">
        Вернуться назад
    </button>
    <hr>

    <h5>%s</h5>
    <hr>

    <table class="table table-bordered">
      %s
      %s
    </table>

  </div>

  <script>
    document.getElementById('backBtn').addEventListener('click', function() {
      window.location.href = '%s';
    });
  </script>
  

</body>
</html>
`,
    title,
    description,
    headers.String(),
    rows.String(),
    redirectBack,
  )

  return []byte(html), nil
}
