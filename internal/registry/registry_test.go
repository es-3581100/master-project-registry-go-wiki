package registry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSHA256String(t *testing.T) {
	if got := SHA256String("abc"); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatal(got)
	}
}
func TestSlugify(t *testing.T) {
	if got := Slugify("Goblin D.M.T."); got != "goblin-d-m-t" {
		t.Fatal(got)
	}
}
func TestStaticExportOmitsHiddenAndProprietary(t *testing.T) {
	tpl, err := Templates()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg := Config{SiteTitle: "T", Owner: "O", StaticAuthEnabled: true}
	ps := []Project{{Slug: "visible", Title: "Visible", Family: "F", Status: "active"}, {Slug: "hidden", Title: "Hidden", Family: "F", Status: "active", HiddenFromDeploy: true}, {Slug: "private", Title: "Private", Family: "F", Status: "active", Proprietary: true}}
	if err := ExportStatic(cfg, ps, tpl, dir, "secret", false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "projects.json"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, "Visible") || strings.Contains(s, "Hidden") || strings.Contains(s, "Private") {
		t.Fatalf("bad public export: %s", s)
	}
}
func TestStaticExportRequiresPassword(t *testing.T) {
	tpl, _ := Templates()
	err := ExportStatic(Config{StaticAuthEnabled: true}, nil, tpl, t.TempDir(), "", false)
	if err == nil {
		t.Fatal("expected password error")
	}
}
