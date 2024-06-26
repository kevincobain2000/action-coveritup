package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/kevincobain2000/action-coveritup/db"
	"github.com/kevincobain2000/action-coveritup/pkg"
)

//go:embed favicon.ico
var favicon embed.FS

//go:embed all:frontend/dist/*
var publicDir embed.FS

//go:embed pkg/migrations/*.sql
var embedMigrations embed.FS

type Flags struct {
	host         string
	port         string
	cors         string
	baseUrl      string
	databaseDSN  string
	databaseName string
	migrate      string
	githubAPI    string
	pprofHost    string
	pprofPort    string
}

var f Flags
var version = "dev"

func main() {
	err := os.Setenv("VERSION", version)
	if err != nil {
		pkg.Logger().Error(err)
	}
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	SetupFlags()
	db.SetupCache()
	db.SetupDatabase()

	if f.migrate != "" {
		err := db.Migrate(f.migrate, embedMigrations, "pkg/migrations")
		if err != nil {
			pkg.Logger().Error(err)
		}
		return
	}
	pkg.StartEcho(pkg.NewEcho(f.baseUrl, publicDir, favicon, f.cors), f.host, f.port)
}

func SetupFlags() {
	flag.StringVar(&f.host, "host", "localhost", "host to serve")
	flag.StringVar(&f.port, "port", "3003", "port to serve")
	flag.StringVar(&f.cors, "cors", "", "cors port to allow")
	flag.StringVar(&f.baseUrl, "base-url", "/", "base url with slash")
	flag.StringVar(&f.databaseDSN, "db-dsn", "root:@tcp(127.0.0.1:3306)/", "databaseURL url")
	flag.StringVar(&f.databaseName, "db-name", "coveritup", "database name")
	flag.StringVar(&f.migrate, "migrate", "", "migrate up, down or redo")
	flag.StringVar(&f.githubAPI, "github-api", "https://api.github.com", "github api url")
	flag.StringVar(&f.pprofHost, "pprof-host", "", "pprof host")
	flag.StringVar(&f.pprofPort, "pprof-port", "", "pprof port")
	flag.Parse()

	if f.databaseDSN != "" && os.Getenv("DATABASE_DSN") == "" {
		err := os.Setenv("DATABASE_DSN", f.databaseDSN)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
	if f.databaseName != "" && os.Getenv("DATABASE_NAME") == "" {
		err := os.Setenv("DATABASE_NAME", f.databaseName)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}

	if f.githubAPI != "" && os.Getenv("GITHUB_API") == "" {
		err := os.Setenv("GITHUB_API", f.githubAPI)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
	if f.baseUrl != "" && os.Getenv("BASE_URL") == "" {
		err := os.Setenv("BASE_URL", f.baseUrl)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}

	if f.pprofHost != "" && os.Getenv("PPROF_HOST") == "" {
		err := os.Setenv("PPROF_HOST", f.pprofHost)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
	if f.pprofPort != "" && os.Getenv("PPROF_PORT") == "" {
		err := os.Setenv("PPROF_PORT", f.pprofPort)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
}
