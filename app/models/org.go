package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/kevincobain2000/action-coveritup/db"
)

type Org struct {
	ID        int64      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"column:name;NOT NULL;size:255" json:"name"`
	CreatedAt *time.Time `gorm:"column;created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
}

func (Org) TableName() string {
	return "orgs"
}

func (o *Org) Get(name string) (*Org, error) {
	var ret Org

	query := `SELECT * FROM orgs WHERE name = @name LIMIT 1`
	err := db.Db().Raw(
		query,
		sql.Named("name", name)).
		Scan(&ret).Error

	ret.Name = strings.TrimSpace(ret.Name)

	return &ret, err
}

func (o *Org) Create(name string) (*Org, error) {
	var ret Org
	name = strings.TrimSpace(name)

	insertQ := `INSERT INTO orgs (name) VALUES (@name)`
	err := db.Db().Raw(
		insertQ,
		sql.Named("name", name)).
		Scan(&ret).Error

	if err != nil {
		return &ret, err
	}

	selectQ := `SELECT * FROM orgs WHERE name = @name LIMIT 1`
	err = db.Db().Raw(
		selectQ,
		sql.Named("name", name)).
		Scan(&ret).Error

	return &ret, err
}
