package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/kevincobain2000/action-coveritup/db"
)

const (
	SAFE_LIMIT_BRANCHES = 100
)

type Coverage struct {
	ID         int64      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	OrgID      int64      `gorm:"column:org_id;NOT NULL" json:"org_id"`
	RepoID     int64      `gorm:"column:repo_id;NOT NULL" json:"repo_id"`
	UserID     int64      `gorm:"column:user_id;NOT NULL" json:"user_id"`
	TypeID     int64      `gorm:"column:type_id;NOT NULL" json:"type_id"`
	PrNum      int64      `gorm:"column:pr_num;NOT NULL" json:"pr_num"`
	BranchName string     `gorm:"column:branch_name;NOT NULL;size:255" json:"branch_name"`
	Commit     string     `gorm:"column:commit;NOT NULL;size:255" json:"commit"`
	Score      float32    `gorm:"column:score;NOT NULL" json:"score"`
	CreatedAt  *time.Time `gorm:"column;created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	DeletedAt  *time.Time `gorm:"column;deleted_at;type:timestamp;" json:"deleted_at"`
}

func (Coverage) TableName() string {
	return "coverages"
}

func (c *Coverage) GetByBranchName(branchName string) ([]Coverage, error) {
	var ret []Coverage
	err := db.Db().Raw("SELECT * FROM coverages WHERE branch_name = @branchName", sql.Named("branchName", branchName)).Scan(&ret).Error
	return ret, err
}

func (c *Coverage) GetAllBranches(orgName string, repoName string, typeName string) ([]string, error) {
	var ret []string
	query := `
	SELECT
		DISTINCT c.branch_name
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	AND
		t.name = @typeName
	ORDER BY FIELD(c.branch_name, 'develop', 'master', 'main') DESC
	LIMIT @limit;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("typeName", typeName),
		sql.Named("limit", SAFE_LIMIT_BRANCHES)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}

	return ret, err

}

type LatestBranchScore struct {
	BranchName string  `json:"branch_name"`
	Commit     string  `json:"commit"`
	Score      float64 `json:"score"`
	TypeName   string  `json:"type_name"`
	Metric     string  `json:"metric"`
	CreatedAt  string  `json:"created_at"`
}

func (c *Coverage) GetLatestBranchScore(orgName string, repoName string, branchName string, typeName string) (LatestBranchScore, error) {

	var ret LatestBranchScore
	query := `
	SELECT
		c.branch_name,
		c.commit,
		Round(c.score, 1) as score,
		t.name as type_name,
		t.metric,
		DATE_FORMAT(c.created_at, '%Y/%m/%d') as created_at
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	AND
		c.branch_name = @branchName
	AND
		t.name = @typeName
	ORDER BY
		c.created_at DESC, c.id DESC
	LIMIT 1;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("branchName", branchName),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	ret.Metric = strings.TrimSpace(ret.Metric)

	return ret, err
}

type LatestBranchScorePR struct {
	BranchName string  `json:"branch_name"`
	Score      float64 `json:"score"`
	TypeName   string  `json:"type_name"`
	Metric     string  `json:"metric"`
	PRNum      int64   `json:"pr_num"`
	CreatedAt  string  `json:"created_at"`
}

func (c *Coverage) GetLatestBranchScores(orgName string, repoName string, branchName string, typeName string) ([]LatestBranchScore, error) {
	var ret []LatestBranchScore
	query := `
	SELECT
		c.branch_name,
		DATE_FORMAT(MAX(c.created_at), '%Y/%m/%d') as created_at,
		Round(c.score, 1) as score
	FROM coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		o.name = @orgName
	AND
		r.name = @repoName
	AND
		c.branch_name = @branchName
	AND
		t.name = @typeName
	GROUP BY
		c.branch_name, DATE_FORMAT(c.created_at, '%Y/%m/%d')
	ORDER BY
		MAX(c.created_at) DESC
	LIMIT 150;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("branchName", branchName),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	return ret, err
}

type LatestPRScoreForCommits struct {
	Commit     string  `json:"commit"`
	BranchName string  `json:"branch_name"`
	Score      float64 `json:"score"`
	TypeName   string  `json:"type_name"`
	Metric     string  `json:"metric"`
}

func (c *Coverage) GetLatestPRScoresForCommits(orgName string, repoName string, prNum int, typeName string) ([]LatestPRScoreForCommits, error) {
	var ret []LatestPRScoreForCommits
	query := `
	SELECT
		c.commit, c.branch_name,
		Round(c.score, 1) as score
	FROM coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		o.name = @orgName
	AND
		r.name = @repoName
	AND
		c.pr_num = @prNum
	AND
		t.name = @typeName
	LIMIT 150;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("prNum", prNum),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	return ret, err
}

type LatestPRScoreForUsers struct {
	ID         int64   `json:"id"`
	UserName   string  `json:"user_name"`
	BranchName string  `json:"branch_name"`
	Score      float64 `json:"score"`
	TypeName   string  `json:"type_name"`
	Metric     string  `json:"metric"`
}

func (c *Coverage) GetLatestPRScoresForUsers(orgName string, repoName string, prNum int, typeName string) ([]LatestPRScoreForUsers, error) {
	var ret []LatestPRScoreForUsers
	query := `
	SELECT
	    c.id,
		u.name as user_name,
		c.branch_name,
		Round(c.score, 1) as score
	FROM coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	LEFT JOIN
		users u ON c.user_id = u.id
	WHERE
		o.name = @orgName
	AND
		r.name = @repoName
	AND
		c.pr_num = @prNum
	AND
		t.name = @typeName
	GROUP BY user_name, branch_name
	ORDER BY c.id ASC
	LIMIT 20;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("prNum", prNum),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	return ret, err
}

type LatestUserScore struct {
	UserName  string  `json:"user_name"`
	Score     float64 `json:"score"`
	TypeName  string  `json:"type_name"`
	Metric    string  `json:"metric"`
	CreatedAt string  `json:"created_at"`
}

func (c *Coverage) GetAllUsers(orgName string, repoName string, typeName string) ([]string, error) {
	var ret []string
	query := `
	SELECT
		DISTINCT u.name as user_name
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		users u ON c.user_id = u.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	AND
		t.name = @typeName
	LIMIT @limit;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("typeName", typeName),
		sql.Named("limit", SAFE_LIMIT_USERS)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}

	return ret, err
}

func (c *Coverage) GetLatestUserScore(orgName string, repoName string, userName string, typeName string) (LatestUserScore, error) {
	var ret LatestUserScore
	query := `
	SELECT
		u.name as user_name,
		Round(c.score, 1) as score,
		t.name as type_name,
		t.metric,
		DATE_FORMAT(c.created_at, '%Y/%m/%d') as created_at
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		users u ON c.user_id = u.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	AND
		u.name = @userName
	AND
		t.name = @typeName
	ORDER BY
		c.created_at DESC, c.id DESC
	LIMIT 1;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("userName", userName),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	ret.Metric = strings.TrimSpace(ret.Metric)

	return ret, err
}

func (c *Coverage) GetLatestUserScores(orgName string, repoName string, userName string, typeName string) ([]LatestUserScore, error) {
	var ret []LatestUserScore
	query := `
	SELECT
		u.name as user_name,
		DATE_FORMAT(MAX(c.created_at), '%Y/%m/%d') as created_at,
		Round(c.score, 1) as score
	FROM coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	LEFT JOIN
		users u ON c.user_id = u.id
	WHERE
		o.name = @orgName
	AND
		r.name = @repoName
	AND
		u.name = @userName
	AND
		t.name = @typeName
	GROUP BY
		c.branch_name, DATE_FORMAT(c.created_at, '%Y/%m/%d')
	ORDER BY
		MAX(c.created_at) DESC
	LIMIT 150;
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("userName", userName),
		sql.Named("typeName", typeName)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}

	return ret, err
}

func (c *Coverage) IsFirstPR(org string, repo string, prNum int) bool {

	var ret = []struct {
		TypeID int `json:"type_id"`
		Total  int `json:"total"`
	}{}
	query := `
	SELECT
		type_id,
		COUNT(c.id) as total
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	WHERE
		o.name = @orgName
	AND
		r.name = @repoName
	AND
		pr_num = @prNum
	GROUP BY
		type_id
	`
	err := db.Db().Raw(
		query,
		sql.Named("orgName", org),
		sql.Named("repoName", repo),
		sql.Named("prNum", prNum)).
		Scan(&ret).Error

	if err != nil {
		return false
	}

	for _, r := range ret {
		if r.Total > 1 {
			return false
		}
	}
	return true

}
func (c *Coverage) DeleteCoveragesByType(org string, repo string, typeName string) error {
	org = strings.TrimSpace(org)
	repo = strings.TrimSpace(repo)
	typeName = strings.TrimSpace(typeName)

	query := `
	DELETE c
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	LEFT JOIN
		types t ON c.type_id = t.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	AND
		t.name = @typeName
	`
	err := db.Db().Exec(
		query,
		sql.Named("orgName", org),
		sql.Named("repoName", repo),
		sql.Named("typeName", typeName)).
		Error
	return err
}

func (c *Coverage) DeleteCoverages(org string, repo string) error {
	org = strings.TrimSpace(org)
	repo = strings.TrimSpace(repo)
	query := `
	DELETE c
	FROM
		coverages c
	LEFT JOIN
		orgs o ON c.org_id = o.id
	LEFT JOIN
		repos r ON c.repo_id = r.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	`
	err := db.Db().Exec(
		query,
		sql.Named("orgName", org),
		sql.Named("repoName", repo)).
		Error

	if err != nil {
		return err
	}
	query = `
	DELETE r
	FROM
		repos r
	LEFT JOIN
		orgs o ON r.org_id = o.id
	WHERE
		r.name = @repoName
	AND
		o.name = @orgName
	`
	err = db.Db().Exec(
		query,
		sql.Named("orgName", org),
		sql.Named("repoName", repo)).
		Error

	if err != nil {
		return err
	}
	query = `
	DELETE o
	FROM
		orgs o
	WHERE
		o.name = @orgName
		AND
		(SELECT COUNT(*) FROM repos r WHERE r.org_id = o.id) = 0
	`
	err = db.Db().Exec(
		query,
		sql.Named("orgName", org)).
		Error

	return err
}

func (c *Coverage) Create(
	orgID int64,
	repoID int64,
	userID int64,
	typeID int64,
	branchName string,
	prNum int64,
	commit string,
	score float32) (*Coverage, error) {
	var ret Coverage

	branchName = strings.TrimSpace(branchName)
	commit = strings.TrimSpace(commit)

	insertQ := `INSERT INTO
		coverages (
			org_id,
			repo_id,
			user_id,
			type_id,
			pr_num,
			commit,
			branch_name,
			score
		)
		VALUES (
			@org_id,
			@repo_id,
			@user_id,
			@type_id,
			@pr_num,
			@commit,
			@branch_name,
			@score
	)`
	err := db.Db().Raw(
		insertQ,
		sql.Named("org_id", orgID),
		sql.Named("repo_id", repoID),
		sql.Named("user_id", userID),
		sql.Named("type_id", typeID),
		sql.Named("pr_num", prNum),
		sql.Named("commit", commit),
		sql.Named("branch_name", branchName),
		sql.Named("score", score)).
		Scan(&ret).Error

	if err != nil {
		return &ret, err
	}

	return &ret, err
}
func (c *Coverage) SoftDeleteCoverages(orgID int64, repoID int64, branches string) error {
	// Split the branches string into a slice
	branchList := strings.Split(branches, " ")

	// Build the NOT IN clause with placeholders
	placeholders := make([]string, len(branchList))
	args := make([]interface{}, len(branchList)+2)

	// Set the orgID and repoID at the beginning of the args slice
	args[0] = orgID
	args[1] = repoID

	// Assign placeholders and arguments for branch names
	for i, branch := range branchList {
		placeholders[i] = "?" // Use "?" as the placeholder for each branch
		args[i+2] = branch    // +2 to account for orgID and repoID at the beginning
	}
	notInClause := strings.Join(placeholders, ", ")

	query := `
	UPDATE
		coverages
	SET
		deleted_at = NOW()
	WHERE
		org_id = ?
	AND
		repo_id = ?
	AND
		branch_name NOT IN (` + notInClause + `)
	`

	// Execute the query
	err := db.Db().Exec(query, args...).Error
	return err
}
