package entity

type Member struct {
	UserId   string
	Username string
	IsActive bool
}

type Team struct {
	TeamName string
	Members  []Member
}
