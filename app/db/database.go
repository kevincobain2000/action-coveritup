package db

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbOnce sync.Once
	db     *gorm.DB
)

func SetupDatabase() {
	syncDb()
}

func Db() *gorm.DB {
	return db
}

func syncDb() *gorm.DB {
	if db != nil {
		return db
	}
	dbOnce.Do(func() {
		dsn := os.Getenv("DATABASE_DSN")
		dbn := os.Getenv("DATABASE_NAME")
		fdsn := dsn + dbn + "?charset=utf8mb4&parseTime=True&loc=Local"

		var err error
		db, err = gorm.Open(mysql.Open(fdsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(1)),
		})
		if err != nil {
			log.Fatal("cannot connect to database")
		}
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatal("cannot get database intance")
		}
		_, err = sqlDb.Exec("SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));")
		if err != nil {
			log.Fatal("cannot set sql_mode")
		}
		configureSQL(sqlDb)
	})

	return db
}

func configureSQL(sqlDB *sql.DB) {
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(30)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(30)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// Default DB close on mysql is 8 hours, so we set way before that (1 min)
	// This can be increased to 1 hour as well
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(1))
}

func Migrate(command string, embedMigrations embed.FS, dir string) error {
	if command == "create" {
		create()
		return nil
	}
	if command == "drop" {
		drop()
		return nil
	}
	driver, err := Db().DB()

	// donot close the connection if you want to add it to part of the app or tests
	// the connection is closed on cli exit anyways
	// defer func() {
	// 	err := driver.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = goose.SetDialect("mysql")
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	err = goose.RunContext(ctx, command, driver, dir)
	if err != nil {
		panic(err)
	}
	return nil
}

func create() {
	dsn := os.Getenv("DATABASE_DSN")
	dbn := os.Getenv("DATABASE_NAME")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbn)
	if err != nil {
		panic(err)
	}
}

func drop() {
	dsn := os.Getenv("DATABASE_DSN")
	dbn := os.Getenv("DATABASE_NAME")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbn)
	if err != nil {
		panic(err)
	}
}
