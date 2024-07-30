package models

import "github.com/google/uuid"

type Message struct {
	Msg string ` json:"msg" db:"message"`
}

type MessageForRepo struct {
	Id  uuid.UUID ` db:"id"`
	Msg string    `db:"message"`
}

type Response struct {
	ProcessedCount    *int32 ` db:"processed_count" json:"processed_count"`
	NotProcessedCount *int32 ` db:"notprocessed_count" json:"notProcessed_count"`
	TotalCount        *int32 ` db:"total_count" json:"total_count"`
}

func (m Message) ConvertToRepoStruct() *MessageForRepo {

	return &MessageForRepo{
		Id:  uuid.New(),
		Msg: m.Msg,
	}

}
