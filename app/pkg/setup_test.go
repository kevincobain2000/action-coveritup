package pkg

import (
	"embed"
	"os"

	"github.com/kevincobain2000/action-coveritup/db"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func BeforeEach() {
	os.Setenv("DATABASE_DSN", "root:@tcp(127.0.0.1:3306)/")
	os.Setenv("DATABASE_NAME", "coverituptest")
	db.Migrate("create", embedMigrations)

	db.SetupDatabase()
	db.Migrate("up", embedMigrations)
}
func AfterEach() {
	os.Setenv("DATABASE_DSN", "root:@tcp(127.0.0.1:3306)/")
	os.Setenv("DATABASE_NAME", "coverituptest")
	db.Migrate("drop", embedMigrations)
}
