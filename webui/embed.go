package webui

import (
	"embed"
	"io/fs"
)

// FS is the single embedded source of truth for registry templates and assets.
//
//go:embed templates/*.html assets/*
var FS embed.FS

func Assets() (fs.FS, error) { return fs.Sub(FS, "assets") }
