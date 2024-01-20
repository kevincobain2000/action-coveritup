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
	baseUrl      string
	databaseDSN  string
	databaseName string
	migrate      string
	githubAPI    string
}

var f Flags
var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	SetupFlags()
	db.SetupCache()
	db.SetupDatabase()

	if f.migrate != "" {
		err := db.Migrate(f.migrate, embedMigrations)
		if err != nil {
			pkg.Logger().Error(err)
		}
		return
	}

	pkg.GracefulServerWithPid(pkg.NewEcho(f.baseUrl, publicDir, favicon), f.host, f.port)
}

func SetupFlags() {
	flag.StringVar(&f.host, "host", "localhost", "host to serve")
	flag.StringVar(&f.port, "port", "3002", "port to serve")
	flag.StringVar(&f.baseUrl, "base-url", "/", "base url with slash")
	flag.StringVar(&f.databaseDSN, "db-dsn", "root:@tcp(127.0.0.1:3306)/", "databaseURL url")
	flag.StringVar(&f.databaseName, "db-name", "coveritup", "database name")
	flag.StringVar(&f.migrate, "migrate", "", "migrate up, down or redo")
	flag.StringVar(&f.githubAPI, "github-api", "https://api.github.com", "github api url")
	flag.Parse()

	if f.databaseDSN != "" {
		err := os.Setenv("DATABASE_DSN", f.databaseDSN)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
	if f.databaseName != "" {
		err := os.Setenv("DATABASE_NAME", f.databaseName)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
	if f.githubAPI != "" {
		err := os.Setenv("GITHUB_API", f.githubAPI)
		if err != nil {
			pkg.Logger().Error(err)
		}
	}
}
