package state

import (
	"context"

	"github.com/google/uuid"
	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/kafka"
	"github.com/rlapenok/messagio/internal/models"
	"github.com/rlapenok/messagio/internal/repo"
)

var State state

func init() {

	cfg := config.New()
	channel := make(chan uuid.UUID, 1000)
	producer, consumer := kafka.CreateProducerAndConsumer(cfg.KafkaConfig, channel)
	consumer.ConsumeMessage()
	State = state{
		repo:     repo.New(cfg.DataBaseConfig, channel),
		producer: producer,
		consumer: consumer,
	}
}

type state struct {
	repo     repo.Repository
	producer kafka.Producer
	consumer kafka.Consumer
}

func (s state) SaveMessage(ctx context.Context, msg models.Message) error {
	//convert sctruct for Repo entity
	msgForRepo := msg.ConvertToRepoStruct()
	//save msg in Repo
	if err := s.repo.Save(ctx, msgForRepo); err != nil {
		return err
	}
	//send msg to kafka
	s.producer.WriteMessage(msgForRepo.Id)
	return nil

}
func (s state) GetSats(ctx context.Context) (*models.Response, error) {
	return s.repo.GetSats(ctx)
}
func (s state) Close() error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	if err := s.consumer.Close(); err != nil {
		return err
	}
	if err := s.repo.Close(); err != nil {
		return err
	}
	return nil

}
