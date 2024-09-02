package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sweatsparks/internal/models"
)

type MatchRepository interface {
	CreateMatch(ctx context.Context, tx *sql.Tx, match *models.Match) error
	FindMatchByUserID(ctx context.Context, tx *sql.Tx, userID1, userID2 uint64) (*models.Match, error)
	FindAllMatchByUserID(ctx context.Context, tx *sql.Tx, userID uint64) ([]*models.Match, error)
}

type MatchRepositoryImpl struct {
}

func NewMatchRepository() MatchRepository {
	return &MatchRepositoryImpl{}
}

func (repository *MatchRepositoryImpl) CreateMatch(ctx context.Context, tx *sql.Tx, match *models.Match) error {
	SQL := `INSERT INTO matches (user_one_id, user_two_id, matched_at) VALUES (?,?,?)`
	response, err := tx.ExecContext(ctx, SQL, match.UserOne, match.UserTwo, match.MatchedTime)
	if err != nil {
		return errors.New("Failed to create a match, transaction rolled back. Reason: " + err.Error())
	}
	matchID, err := response.LastInsertId()
	if err != nil {
		return errors.New("Failed to retrieve match_id, transaction rolled back. Reason:" + err.Error())
	}

	match.Id = uint64(matchID)
	return nil
}

func (repository *MatchRepositoryImpl) FindMatchByUserID(ctx context.Context, tx *sql.Tx, userID1, userID2 uint64) (*models.Match, error) {
	SQL := "select id, user_one_id, user_two_id, matched_at from matches where (user_one_id = ? and user_two_id = ?) or (user_one_id = ? and user_two_id = ?)"
	rows, err := tx.QueryContext(ctx, SQL, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var match = models.Match{}
	if rows.Next() {
		err := rows.Scan(&match.Id, &match.UserOne, &match.UserTwo, &match.MatchedTime)
		if err != nil {
			return nil, err
		}
		return &match, nil
	} else {
		return nil, errors.New("match is not found")
	}
}
func (repository *MatchRepositoryImpl) FindAllMatchByUserID(ctx context.Context, tx *sql.Tx, userID uint64) ([]*models.Match, error) {
	SQL := "select id, user_one_id, user_two_id, matched_at from matches where (user_one_id = ? or user_two_id = ?)"
	rows, err := tx.QueryContext(ctx, SQL, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*models.Match
	for rows.Next() {
		var match models.Match
		if err := rows.Scan(&match.Id, &match.UserOne, &match.UserTwo, &match.MatchedTime); err != nil {
			return nil, err
		}
		matches = append(matches, &match)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}
