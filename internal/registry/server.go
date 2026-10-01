package registry

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Server struct {
	cfg   Config
	store *Store
	tpl   *template.Template
}

func NewServer(cfg Config, store *Store, tpl *template.Template) *Server {
	return &Server{cfg: cfg, store: store, tpl: tpl}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	assets, _ := Assets()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("/unlock", s.handleUnlock)
	mux.HandleFunc("/project/", s.auth(s.handleProject))
	mux.HandleFunc("/admin/edit", s.auth(s.handleAdminEdit))
	mux.HandleFunc("/admin/save", s.auth(s.handleAdminSave))
	mux.HandleFunc("/admin/retire", s.auth(s.handleAdminRetire))
	mux.HandleFunc("/admin/toggle-deploy", s.auth(s.handleAdminToggleDeploy))
	mux.HandleFunc("/admin", s.auth(s.handleAdmin))
	mux.HandleFunc("/api/projects.json", s.auth(s.handleAPI))
	mux.HandleFunc("/", s.auth(s.handleIndex))
	return securityHeaders(mux)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self'; script-src 'self'; form-action 'self'; base-uri 'self'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.cfg.LocalAuthEnabled {
			next(w, r)
			return
		}
		c, err := r.Cookie("mpr_auth")
		if err == nil && s.cfg.LocalPasswordSHA256 != "" && subtle.ConstantTimeCompare([]byte(c.Value), []byte(s.cfg.LocalPasswordSHA256)) == 1 {
			next(w, r)
			return
		}
		http.Redirect(w, r, "/unlock?next="+url.QueryEscape(r.URL.RequestURI()), http.StatusSeeOther)
	}
}

func (s *Server) handleUnlock(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if s.cfg.LocalPasswordSHA256 == "" {
			http.Error(w, "local auth is enabled but localPasswordSHA256 is empty", 500)
			return
		}
		if subtle.ConstantTimeCompare([]byte(SHA256String(r.FormValue("password"))), []byte(s.cfg.LocalPasswordSHA256)) == 1 {
			http.SetCookie(w, &http.Cookie{Name: "mpr_auth", Value: s.cfg.LocalPasswordSHA256, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Expires: time.Now().Add(12 * time.Hour)})
			dest := r.FormValue("next")
			if dest == "" {
				dest = "/"
			}
			http.Redirect(w, r, dest, http.StatusSeeOther)
			return
		}
	}
	next := r.URL.Query().Get("next")
	fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8"><link rel="stylesheet" href="/assets/style.css"><title>Unlock</title></head><body><main class="unlock"><section class="glass"><div class="eyebrow">LOCAL GATE</div><h1>Master Project Registry</h1><form method="post"><input type="hidden" name="next" value="%s"><label>Password<input type="password" name="password" autofocus></label><button>Unlock</button></form></section></main></body></html>`, template.HTMLEscapeString(next))
}

func (s *Server) viewData(projects []Project) ViewData {
	return ViewData{Config: s.cfg, Projects: projects, Families: Families(projects), Statuses: Statuses(projects)}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := s.viewData(s.store.All())
	if err := s.tpl.ExecuteTemplate(w, "index", data); err != nil {
		log.Println(err)
	}
}
func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/project/")
	p, ok := s.store.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := s.viewData(s.store.All())
	data.Project = &p
	if err := s.tpl.ExecuteTemplate(w, "project", data); err != nil {
		log.Println(err)
	}
}
func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.store.All())
}
func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	data := s.viewData(s.store.All())
	data.Admin = true
	if err := s.tpl.ExecuteTemplate(w, "admin", data); err != nil {
		log.Println(err)
	}
}
func (s *Server) handleAdminEdit(w http.ResponseWriter, r *http.Request) {
	var p Project
	if slug := r.URL.Query().Get("slug"); slug != "" {
		var ok bool
		p, ok = s.store.Get(slug)
		if !ok {
			http.NotFound(w, r)
			return
		}
	}
	data := s.viewData(s.store.All())
	data.Admin = true
	data.Project = &p
	if err := s.tpl.ExecuteTemplate(w, "admin-edit", data); err != nil {
		log.Println(err)
	}
}
func parseBool(v string) bool { return v == "on" || v == "true" || v == "1" }
func splitTags(s string) []string {
	var out []string
	for _, x := range strings.Split(s, ",") {
		x = strings.TrimSpace(x)
		if x != "" {
			out = append(out, x)
		}
	}
	return out
}
func (s *Server) handleAdminSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", 405)
		return
	}
	_ = r.ParseForm()
	p := Project{ID: r.FormValue("id"), Slug: r.FormValue("slug"), Title: r.FormValue("title"), Family: r.FormValue("family"), Status: r.FormValue("status"), Kind: r.FormValue("kind"), Summary: r.FormValue("summary"), Path: r.FormValue("path"), RepoURL: r.FormValue("repoURL"), SiteURL: r.FormValue("siteURL"), Tags: splitTags(r.FormValue("tags")), Notes: r.FormValue("notes"), Source: r.FormValue("source"), HiddenFromDeploy: parseBool(r.FormValue("hiddenFromDeploy")), Proprietary: parseBool(r.FormValue("proprietary"))}
	if err := s.store.Upsert(p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
func (s *Server) handleAdminRetire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", 405)
		return
	}
	slug := r.FormValue("slug")
	_ = s.store.Mutate(slug, func(p *Project) { p.Status = "retired" })
	http.Redirect(w, r, "/admin", 303)
}
func (s *Server) handleAdminToggleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", 405)
		return
	}
	slug := r.FormValue("slug")
	_ = s.store.Mutate(slug, func(p *Project) { p.HiddenFromDeploy = !p.HiddenFromDeploy })
	http.Redirect(w, r, "/admin", 303)
}

func Copy(dst io.Writer, src io.Reader) error { _, err := io.Copy(dst, src); return err }
