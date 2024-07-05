package forms

import (
  "bytes"
  _ "embed"
  "fmt"
  "html/template"

  "gopkg.in/yaml.v3"
)

var (
  prognosisBuf []byte

  //go:embed prognosis.template
  prognosisTemplateString string

  //go:embed prognosis.yaml
  prognosisConfigBuf []byte
)

type prognosisConfig struct {
  Sections []prognosisSection `yaml:"sections"`
}

type prognosisSection struct {
  Questions []prognosisQuestion `yaml:"questions"`
}

type prognosisQuestion struct {
  Label   string            `yaml:"label"`
  Choises []prognosisChoise `yaml:"choises"`
}

type prognosisChoise struct {
  Label string  `yaml:"label"`
  Value float64 `yaml:"value"`
}

func buildPrognosis() ([]byte, error) {
  config := new(prognosisConfig)

  if err := yaml.Unmarshal(prognosisConfigBuf, config); err != nil {
    return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
  }

  p, err := template.New("prognosis").
    Funcs(funcMap).
    Parse(prognosisTemplateString)

  if err != nil {
    return nil, fmt.Errorf("template.New.Parse: %w", err)
  }

  buf := new(bytes.Buffer)

  if err = p.Execute(buf, config); err != nil {
    return nil, fmt.Errorf("parsed.Execute: %w", err)
  }

  return buf.Bytes(), nil
}
