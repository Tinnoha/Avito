package entity

type Member struct {
	UserId   string
	Username string
	IsActive bool
}

type Teams struct {
	TeamName string
	Members  []Member
}
