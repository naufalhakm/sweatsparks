package repositories

import (
	"context"
	"database/sql"
	"sweatsparks/internal/models"
)

type MessageRepository interface {
	GetMessageByMatchID(ctx context.Context, tx *sql.Tx, matchID int) ([]*models.Message, error)
}

type MessageRepositoryImpl struct{}

func NewMessageRepository() MessageRepository {
	return &MessageRepositoryImpl{}
}

func (repository *MessageRepositoryImpl) GetMessageByMatchID(ctx context.Context, tx *sql.Tx, matchID int) ([]*models.Message, error) {
	SQL := `SELECT id, match_id, sender_id, content, sent_at FROM messages WHERE match_id = ? ORDER BY sent_at ASC`

	rows, err := tx.QueryContext(ctx, SQL, matchID)
	if err != nil {
		return nil, err
	}

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(
			&message.Id,
			&message.MatchID,
			&message.SendAt,
			&message.Content,
			&message.SendAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}
	return messages, nil
}
