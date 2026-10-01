package registry

import (
	"html/template"
	"io/fs"
	"strings"

	"github.com/es-3581100/master-project-registry-go-wiki/webui"
)

type ViewData struct {
	Config     Config
	Project    *Project
	Projects   []Project
	Families   []string
	Statuses   []string
	Admin      bool
	Static     bool
	StaticGate bool
}

func Templates() (*template.Template, error) {
	fm := template.FuncMap{
		"join":        func(v []string, sep string) string { return strings.Join(v, sep) },
		"statusClass": func(s string) string { return "status-" + Slugify(s) },
		"splitLines":  func(s string) []string { return strings.Split(s, "\n") },
	}
	return template.New("root").Funcs(fm).ParseFS(webui.FS, "templates/*.html")
}

func Assets() (fs.FS, error) { return webui.Assets() }
