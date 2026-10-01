package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/es-3581100/master-project-registry-go-wiki/internal/registry"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "serve":
		serve(os.Args[2:])
	case "export":
		export(os.Args[2:])
	case "hash":
		hash(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}
func usage() { fmt.Println("registry serve|export|hash") }
func common(fs *flag.FlagSet) (*string, *string) {
	c := fs.String("config", "config.json", "config path")
	d := fs.String("data", "data/projects.json", "projects data path")
	return c, d
}
func serve(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	cp, dp := common(fs)
	bind := fs.String("bind", "", "override bind address")
	_ = fs.Parse(args)
	cfg, err := registry.LoadConfig(*cp)
	if err != nil {
		log.Fatal(err)
	}
	if *bind != "" {
		cfg.Bind = *bind
	}
	store, err := registry.OpenStore(*dp)
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := registry.Templates()
	if err != nil {
		log.Fatal(err)
	}
	srv := registry.NewServer(cfg, store, tpl)
	log.Printf("Master Project Registry local wiki: http://%s", cfg.Bind)
	log.Printf("Admin: http://%s/admin", cfg.Bind)
	log.Fatal(http.ListenAndServe(cfg.Bind, srv.Handler()))
}
func export(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	cp, dp := common(fs)
	out := fs.String("out", "docs", "static output directory")
	password := fs.String("password", "", "static gate password (or REGISTRY_DEPLOY_PASSWORD)")
	noGate := fs.Bool("no-gate", false, "explicitly disable static password gate")
	_ = fs.Parse(args)
	if *password == "" {
		*password = os.Getenv("REGISTRY_DEPLOY_PASSWORD")
	}
	cfg, err := registry.LoadConfig(*cp)
	if err != nil {
		log.Fatal(err)
	}
	store, err := registry.OpenStore(*dp)
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := registry.Templates()
	if err != nil {
		log.Fatal(err)
	}
	abs, _ := filepath.Abs(*out)
	if err := registry.ExportStatic(cfg, store.All(), tpl, *out, *password, *noGate); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("static wiki exported to %s\n", abs)
}
func hash(args []string) {
	fs := flag.NewFlagSet("hash", flag.ExitOnError)
	v := fs.String("value", "", "value to SHA-256")
	_ = fs.Parse(args)
	if *v == "" {
		log.Fatal("--value required")
	}
	fmt.Println(registry.SHA256String(*v))
}
