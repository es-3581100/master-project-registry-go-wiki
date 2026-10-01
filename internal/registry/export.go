package registry

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func ExportStatic(cfg Config, projects []Project, tpl *template.Template, outDir, password string, noGate bool) error {
	if cfg.StaticAuthEnabled && !noGate && password == "" {
		return fmt.Errorf("static password required: pass --password or set REGISTRY_DEPLOY_PASSWORD (or explicitly use --no-gate)")
	}
	pub := make([]Project, 0, len(projects))
	for _, p := range projects {
		if !p.HiddenFromDeploy && !p.Proprietary {
			pub = append(pub, p)
		}
	}
	SortProjects(pub)
	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(outDir, "projects"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(outDir, "assets"), 0o755); err != nil {
		return err
	}
	assets, _ := Assets()
	if err := fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(assets, path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(outDir, "assets", filepath.Base(path)), b, 0o644)
	}); err != nil {
		return err
	}
	_ = os.WriteFile(filepath.Join(outDir, ".nojekyll"), []byte{}, 0o644)
	sc := StaticConfig{SiteTitle: cfg.SiteTitle, GateEnabled: cfg.StaticAuthEnabled && !noGate}
	if sc.GateEnabled {
		sc.PasswordSHA256 = SHA256String(password)
	}
	b, _ := json.MarshalIndent(sc, "", "  ")
	_ = os.WriteFile(filepath.Join(outDir, "site-config.json"), append(b, '\n'), 0o644)
	pb, _ := json.MarshalIndent(pub, "", "  ")
	_ = os.WriteFile(filepath.Join(outDir, "projects.json"), append(pb, '\n'), 0o644)
	vd := ViewData{Config: cfg, Projects: pub, Families: Families(pub), Statuses: Statuses(pub), Static: true, StaticGate: sc.GateEnabled}
	if err := writeTemplate(filepath.Join(outDir, "index.html"), tpl, "index", vd); err != nil {
		return err
	}
	for i := range pub {
		p := pub[i]
		d := vd
		d.Project = &p
		if err := writeTemplate(filepath.Join(outDir, "projects", p.Slug+".html"), tpl, "project", d); err != nil {
			return err
		}
	}
	return nil
}

func writeTemplate(path string, tpl *template.Template, name string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.ExecuteTemplate(f, name, data)
}
