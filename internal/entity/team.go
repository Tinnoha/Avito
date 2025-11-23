package entity

type Member struct {
	UserId   string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type Team struct {
	TeamName string   `db:"team_name" json:"team_name"`
	Members  []Member `json:"members"`
}
