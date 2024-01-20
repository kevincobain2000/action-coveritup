package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/kevincobain2000/action-coveritup/db"
)

type Repo struct {
	ID        int64      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	OrgID     int64      `gorm:"column:org_id;NOT NULL" json:"org_id"`
	Name      string     `gorm:"column:name;NOT NULL;size:255" json:"name"`
	CreatedAt *time.Time `gorm:"column;created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
}

func (Repo) TableName() string {
	return "repos"
}
func (r *Repo) Get(orgID int64, name string) (Repo, error) {
	var ret Repo

	query := `SELECT * FROM repos WHERE org_id = @org_id AND name = @name LIMIT 1`
	err := db.Db().Raw(
		query,
		sql.Named("org_id", orgID),
		sql.Named("name", name)).
		Scan(&ret).Error

	ret.Name = strings.TrimSpace(ret.Name)

	return ret, err
}

func (r *Repo) Create(orgID int64, name string) (Repo, error) {
	var ret Repo

	query := `INSERT INTO repos (org_id, name) VALUES (@org_id, @name)`
	err := db.Db().Raw(
		query,
		sql.Named("org_id", orgID),
		sql.Named("name", name)).
		Scan(&ret).Error

	return ret, err
}
