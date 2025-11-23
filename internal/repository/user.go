package repository

import (
	"Avito/internal/entity"

	"github.com/jmoiron/sqlx"
	_ "github.com/powerman/pqx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Save(vasya entity.User, id int) (entity.User, error) {
	_, err := u.db.Exec("INSERT INTO users (user_id, username, team_id, is_active) values ($1, $2, $3, $4)", vasya.UserId, vasya.Username, id, true)

	if err != nil {
		return entity.User{}, err
	}

	return vasya, nil
}
func (u *userRepo) SetIsActive(userId string, isActive bool) (entity.User, error) {
	var vasya entity.User
	_, err := u.db.Exec(`UPDATE users
	set is_active = $1
	where user_id = $2
	`, isActive, userId)

	if err != nil {
		return entity.User{}, err
	}

	err = u.db.QueryRow(`select 
	users.user_id,
	users.username,
	users.is_active,
	teams.team_name
	from users
	join teams on users.team_id=teams.id
	where users.user_id=$1
	`, userId).Scan(&vasya.UserId, &vasya.Username, &vasya.IsActive, &vasya.TeamName)

	if err != nil {
		return entity.User{}, err
	}

	return vasya, nil

}
func (u *userRepo) IsExists(userId string) bool {
	var exists bool

	err := u.db.QueryRow(`
        select exists(
            select 1 from users where user_id = $1
        )
    `, userId).Scan(&exists)

	return err == nil && exists
}
func (u *userRepo) UserById(userId string) (entity.User, error) {
	var vasya entity.User

	err := u.db.QueryRow(`select 
	users.user_id,
	users.username,
	users.is_active,
	teams.team_name
	from users
	join teams on users.team_id=teams.id
	where users.user_id=$1
	`, userId).Scan(&vasya.UserId, &vasya.Username, &vasya.IsActive, &vasya.TeamName)

	if err != nil {
		return entity.User{}, err
	}

	return vasya, err
}
