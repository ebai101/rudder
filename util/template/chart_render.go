package template

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/render"
	tpls "github.com/go-echarts/go-echarts/v2/templates"
)

const HeaderTpl = `
{{ define "header" }}
<head>
{{- range .JSAssets.Values }}
  <script src="{{ . }}"></script>
{{- end }}
{{- range .CSSAssets.Values }}
  <link href="{{ . }}" rel="stylesheet">
{{- end }}
</head>
{{ end }}
`

func ChartComponent(chart Renderable) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return chart.Render(w)
	})
}

type chartRender struct {
	render.BaseRender
	c      interface{}
	before []func()
}

func NewChartRender(c interface{}, before ...func()) render.Renderer {
	return &chartRender{c: c, before: before}
}

func (r *chartRender) Render(w io.Writer) error {
	for _, fn := range r.before {
		fn()
	}

	contents := []string{HeaderTpl, tpls.BaseTpl, tpls.ChartTpl}
	tpl := render.MustTemplate("chart", contents)

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "chart", r.c); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}
