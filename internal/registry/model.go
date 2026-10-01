package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"
)

type Project struct {
	ID               string   `json:"id"`
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	Family           string   `json:"family"`
	Status           string   `json:"status"`
	Kind             string   `json:"kind"`
	Summary          string   `json:"summary"`
	Path             string   `json:"path"`
	RepoURL          string   `json:"repoURL"`
	SiteURL          string   `json:"siteURL"`
	Tags             []string `json:"tags"`
	Notes            string   `json:"notes"`
	Source           string   `json:"source"`
	HiddenFromDeploy bool     `json:"hiddenFromDeploy"`
	Proprietary      bool     `json:"proprietary"`
}

type Config struct {
	SiteTitle           string `json:"siteTitle"`
	Owner               string `json:"owner"`
	Bind                string `json:"bind"`
	LocalAuthEnabled    bool   `json:"localAuthEnabled"`
	LocalPasswordSHA256 string `json:"localPasswordSHA256"`
	StaticAuthEnabled   bool   `json:"staticAuthEnabled"`
}

type StaticConfig struct {
	SiteTitle      string `json:"siteTitle"`
	GateEnabled    bool   `json:"gateEnabled"`
	PasswordSHA256 string `json:"passwordSHA256"`
}

func SHA256String(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

var slugRE = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRE.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "project"
	}
	return s
}

func SortProjects(in []Project) {
	sort.SliceStable(in, func(i, j int) bool {
		if in[i].Family != in[j].Family {
			return strings.ToLower(in[i].Family) < strings.ToLower(in[j].Family)
		}
		return strings.ToLower(in[i].Title) < strings.ToLower(in[j].Title)
	})
}

func Families(in []Project) []string {
	m := map[string]bool{}
	for _, p := range in {
		if p.Family != "" {
			m[p.Family] = true
		}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func Statuses(in []Project) []string {
	m := map[string]bool{}
	for _, p := range in {
		if p.Status != "" {
			m[p.Status] = true
		}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
