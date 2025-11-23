package repository

import (
	"Avito/internal/entity"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type PullRequestRepo struct {
	db *sqlx.DB
}

func NewPullRequestRepo(db *sqlx.DB) *PullRequestRepo {
	return &PullRequestRepo{
		db: db,
	}
}

func (t *PullRequestRepo) Create(request entity.ShortPullRequest, reviewers []string) (entity.PullRequest, error) {
	fullRequest := entity.PullRequest{}

	err := t.db.QueryRow(`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id,status, created_at) 
	values ($1, $2, $3,$4, $5)
	returning pull_request_id, pull_request_name, author_id,status, created_at`,
		request.PullRequestId, request.PullRequestName, request.AuthorId, "open", time.Now()).Scan(
		&fullRequest.PullRequestId,
		&fullRequest.PullRequestName,
		&fullRequest.AuthorId,
		&fullRequest.Status,
		&fullRequest.CreatedAt,
	)

	if err != nil {
		return entity.PullRequest{}, err
	}

	fullRequest.MergeAt = nil
	for _, v := range reviewers {
		_, err = t.db.Exec(`insert into pull_reviewer (pull_request_id, reviewer_id) 
		values ($1, $2)`, request.PullRequestId, v)

		if err != nil {
			return entity.PullRequest{}, err
		}
	}
	fullRequest.AssignedReviewers = reviewers
	return fullRequest, nil
}

func (t *PullRequestRepo) Merge(PullRequestId string) (entity.PullRequest, error) {
	fullRequest := entity.PullRequest{}
	_, err := t.db.Exec(`update pull_requests set merged_at = $1 where pull_request_id = $2`, time.Now(), PullRequestId)

	if err != nil {
		return entity.PullRequest{}, err
	}

	err = t.db.QueryRow(`select
	pull_request_id,
	pull_request_name,
	author_id,
	status,
	created_at,
	merged_at
	from pull_requests
	where pull_request_id = $1`, PullRequestId).Scan(
		&fullRequest.PullRequestId,
		&fullRequest.PullRequestName,
		&fullRequest.AuthorId,
		&fullRequest.Status,
		&fullRequest.CreatedAt,
		&fullRequest.MergeAt,
	)

	if err != nil {
		return entity.PullRequest{}, err
	}

	rows, err := t.db.Query(`select 
	reviewer_id
	from pull_reviewer
	where pull_request_id = $1`, PullRequestId)

	if err != nil {
		return entity.PullRequest{}, err
	}

	defer rows.Close()

	reviewers := []string{}

	for rows.Next() {
		var member string
		err := rows.Scan(&member)

		if err != nil {
			return entity.PullRequest{}, err
		}

		reviewers = append(reviewers, member)
	}

	fullRequest.AssignedReviewers = reviewers

	return fullRequest, nil
}
func (t *PullRequestRepo) Reassign(PullRequestId string, oldUserId string, newUserId string) error {
	_, err := t.db.Exec(`update pull_reviewer set reviewer_id = $1 where reviewer_id = $2 and pull_request_id = $3`, newUserId, oldUserId, PullRequestId)

	if err != nil {
		return err
	}

	return nil
}
func (t *PullRequestRepo) RequestsById(UserId string, all bool) ([]entity.ShortPullRequest, error) {
	var rows *sql.Rows
	var err error
	if all {
		rows, err = t.db.Query(`select 
		pull_requests.pull_request_id,
		pull_requests.pull_request_name,
		pull_requests.author_id,
		pull_requests.status
		from pull_requests
		join pull_reviewer on pull_reviewer.pull_request_id=pull_requests.pull_request_id
		where pull_reviewer.reviewer_id = $1 `, UserId)
	} else {
		rows, err = t.db.Query(`select 
			pull_requests.pull_request_id,
			pull_requests.pull_request_name,
			pull_requests.author_id,
			pull_requests.status
			from pull_requests
			join pull_reviewer on pull_reviewer.pull_request_id=pull_requests.pull_request_id
			where pull_reviewer.reviewer_id = $1 and status = 'open' `, UserId)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	requests := []entity.ShortPullRequest{}

	for rows.Next() {
		request := entity.ShortPullRequest{}
		err = rows.Scan(&request.PullRequestId, &request.PullRequestName, &request.AuthorId, &request.Status)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}
func (t *PullRequestRepo) IsExists(PullRequestId string) bool {
	var exists bool

	err := t.db.QueryRow(`
        select exists(
            select 1 from pull_requests where pull_request_id = $1
        )
    `, PullRequestId).Scan(&exists)

	return err == nil && exists
}
func (t *PullRequestRepo) RequestByID(PullRequestId string) (entity.PullRequest, error) {
	fullRequest := entity.PullRequest{}
	err := t.db.QueryRow(`select
	pull_request_id,
	pull_request_name,
	author_id,
	status,
	created_at,
	merged_at
	from pull_requests
	where pull_request_id = $1`, PullRequestId).Scan(
		&fullRequest.PullRequestId,
		&fullRequest.PullRequestName,
		&fullRequest.AuthorId,
		&fullRequest.Status,
		&fullRequest.CreatedAt,
		&fullRequest.MergeAt,
	)

	if err != nil {
		return entity.PullRequest{}, err
	}

	rows, err := t.db.Query(`select 
	reviewer_id
	from pull_reviewer
	where pull_request_id = $1`, PullRequestId)

	if err != nil {
		return entity.PullRequest{}, err
	}

	defer rows.Close()

	reviewers := []string{}

	for rows.Next() {
		var member string
		err := rows.Scan(&member)

		if err != nil {
			return entity.PullRequest{}, err
		}

		reviewers = append(reviewers, member)
	}

	fullRequest.AssignedReviewers = reviewers

	return fullRequest, nil
}
func (t *PullRequestRepo) IsMerged(PullRequestId string) (bool, error) {
	var mergeTime *time.Time

	err := t.db.QueryRow("select merged_at from pull_requests where pull_request_id = $1", PullRequestId).Scan(&mergeTime)
	if err != nil {
		return false, err
	}

	if mergeTime != nil {
		return true, nil
	}
	return false, nil

}
