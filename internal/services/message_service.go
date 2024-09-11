package services

import (
	"context"
	"database/sql"
	"sweatsparks/internal/commons/response"
	"sweatsparks/internal/params"
	"sweatsparks/internal/repositories"
	"sweatsparks/pkg/helpers"
)

type MessageService interface {
	GetMessageByMatchId(ctx context.Context, id int) ([]*params.MessageResponse, *response.CustomError)
}

type MessageServiceImpl struct {
	MySqlDB           *sql.DB
	MessageRepository repositories.MessageRepository
}

func NewMessageService(db *sql.DB, messageRepository repositories.MessageRepository) MessageService {
	return &MessageServiceImpl{
		MySqlDB:           db,
		MessageRepository: messageRepository,
	}
}

func (service *MessageServiceImpl) GetMessageByMatchId(ctx context.Context, id int) ([]*params.MessageResponse, *response.CustomError) {
	tx, err := service.MySqlDB.Begin()
	if err != nil {
		return nil, response.GeneralErrorWithAdditionalInfo("Failed Connection to MySQL Errors: %s", err.Error())
	}
	defer helpers.CommitOrRollback(tx)

	messages, err := service.MessageRepository.GetMessageByMatchID(ctx, tx, id)
	if err != nil {
		return nil, response.BadRequestErrorWithAdditionalInfo("messages not found.")
	}

	var result []*params.MessageResponse
	for _, msg := range messages {
		result = append(result, &params.MessageResponse{
			Id:       msg.Id,
			MatchID:  msg.MatchID,
			SenderID: msg.SenderID,
			Content:  msg.Content,
			SendAt:   msg.SendAt,
		})
	}

	return result, nil
}
