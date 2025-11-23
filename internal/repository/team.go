package repository

import (
	"Avito/internal/entity"

	"github.com/jmoiron/sqlx"
)

type TeamRepo struct {
	db *sqlx.DB
}

func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (t *TeamRepo) GetReviewes(AuthorId string, teamName string) ([]string, error) {
	activeMembers := []string{}
	rows, err := t.db.Query(`SELECT 
	users.user_id
	from users
	join teams on teams.id=users.team_id
	where is_active = true and teams.team_name=$1 and users.user_id != $2`, teamName, AuthorId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var vasya string
		err := rows.Scan(&vasya)

		if err != nil {
			return nil, err
		}

		activeMembers = append(activeMembers, vasya)
	}

	return activeMembers, nil
}
func (t *TeamRepo) NewReviewer(AuthorId string, OldReviewer string, TeamName string, PullRequestId string) ([]string, error) {
	activeMembers := []string{}

	rows, err := t.db.Query(`
        SELECT users.user_id
        FROM users
        JOIN teams ON teams.id = users.team_id
        WHERE users.is_active = true 
          AND teams.team_name = $1 
          AND users.user_id != $2 
          AND users.user_id != $3
          AND users.user_id NOT IN (
              SELECT reviewer_id 
              FROM pull_reviewer 
              WHERE pull_request_id = $4
          )`,
		TeamName, AuthorId, OldReviewer, PullRequestId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var member string
		err := rows.Scan(&member)
		if err != nil {
			return nil, err
		}
		activeMembers = append(activeMembers, member)
	}

	return activeMembers, nil
}
func (t *TeamRepo) Create(team entity.Team) (entity.Team, int, error) {
	var id int
	err := t.db.QueryRow("INSERT INTO teams (team_name) values ($1) returning id", team.TeamName).Scan(&id)
	if err != nil {
		return entity.Team{}, 0, err
	}

	return team, id, nil
}
func (t *TeamRepo) GetByName(TeamName string) (entity.Team, error) {
	var team entity.Team

	err := t.db.QueryRow("select team_name from teams where team_name = $1", TeamName).Scan(&team.TeamName)

	if err != nil {
		return entity.Team{}, err
	}

	members := []entity.Member{}

	rows, err := t.db.Query(`select 
	users.user_id, 
	users.username, 
	users.is_active 
	from users 
	join teams on teams.id=users.team_id 
	where teams.team_name = $1`, TeamName)

	if err != nil {
		return entity.Team{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var vasya entity.Member
		err = rows.Scan(&vasya.UserId, &vasya.Username, &vasya.IsActive)

		if err != nil {
			return entity.Team{}, err
		}

		members = append(members, vasya)
	}

	team.Members = members

	return team, nil

}
func (t *TeamRepo) IsExists(TeamName string) bool {
	var exists bool

	err := t.db.QueryRow(`
        select exists(
            select from teams where team_name = $1
        )
    `, TeamName).Scan(&exists)

	return err == nil && exists
}
