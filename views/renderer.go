package views

import (
	"html/template"
	"io"
	"time"

	"github.com/steemwatch/status/checks"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Context struct {
	Sections []*checks.SectionSummary
}

type Template struct {
	templates *template.Template
}

func NewRenderer(templateGlob string) (*Template, error) {
	funcMap := template.FuncMap{
		"systemOK": func(results []*checks.SectionSummary) bool {
			for _, section := range results {
				for _, check := range section.Checks {
					if check.Result != checks.CheckResultPassing {
						return false
					}
				}
			}
			return true
		},
		"sectionStyle": func(summary *checks.SectionSummary) string {
			var (
				somePending bool
				somePassing bool
				someFailing bool
			)
			for _, check := range summary.Checks {
				switch check.Result {
				case checks.CheckResultPending:
					somePending = true
				case checks.CheckResultPassing:
					somePassing = true
				case checks.CheckResultFailing:
					someFailing = true
				}
			}

			switch {
			case somePending:
				return "default"
			case somePassing && someFailing:
				return "warning"
			case somePassing:
				return "success"
			case someFailing:
				return "danger"
			default:
				return "default"
			}
		},
		"resultToStyle": func(result checks.CheckResult) string {
			switch result {
			case checks.CheckResultPending:
				return "info"
			case checks.CheckResultPassing:
				return "success"
			case checks.CheckResultFailing:
				return "danger"
			default:
				return "default"
			}
		},
		"delta": func(timestamp time.Time) string {
			return time.Now().Sub(timestamp).String()
		},
	}

	t, err := template.New("").Funcs(funcMap).ParseGlob(templateGlob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load view templates")
	}

	return &Template{t}, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
