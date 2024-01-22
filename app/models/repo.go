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
func (r *Repo) Get(orgID int64, name string) (*Repo, error) {
	var ret Repo
	name = strings.TrimSpace(name)

	query := `SELECT * FROM repos WHERE org_id = @org_id AND name = @name LIMIT 1`
	err := db.Db().Raw(
		query,
		sql.Named("org_id", orgID),
		sql.Named("name", name)).
		Scan(&ret).Error

	return &ret, err
}

func (r *Repo) Create(orgID int64, name string) (*Repo, error) {
	var ret Repo
	name = strings.TrimSpace(name)

	insertQ := `INSERT INTO repos (org_id, name) VALUES (@org_id, @name)`
	err := db.Db().Raw(
		insertQ,
		sql.Named("org_id", orgID),
		sql.Named("name", name)).
		Scan(&ret).Error

	if err != nil {
		return &ret, err
	}
	selectQ := `SELECT * FROM repos WHERE org_id = @orgID AND name = @name LIMIT 1`
	err = db.Db().Raw(
		selectQ,
		sql.Named("orgID", orgID),
		sql.Named("name", name)).
		Scan(&ret).Error

	return &ret, err
}
