package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/kevincobain2000/action-coveritup/db"
)

const (
	SAFE_LIMIT_USERS = 100
)

type User struct {
	ID        int64      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name      string     `gorm:"column:name;NOT NULL;size:255" json:"name"`
	CreatedAt *time.Time `gorm:"column;created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Get(name string) (User, error) {
	var ret User

	query := `SELECT * FROM users WHERE name = @name LIMIT 1`
	err := db.Db().Raw(
		query,
		sql.Named("name", name)).
		Scan(&ret).Error

	ret.Name = strings.TrimSpace(ret.Name)

	return ret, err
}

func (u *User) Create(name string) (User, error) {
	var ret User

	query := `INSERT INTO users (name) VALUES (@name)`
	err := db.Db().Raw(
		query,
		sql.Named("name", name)).
		Scan(&ret).Error

	return ret, err
}
