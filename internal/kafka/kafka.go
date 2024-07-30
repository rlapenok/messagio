package kafka

import (
	"github.com/google/uuid"
	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/utils"
)

func CreateProducerAndConsumer(cfg config.KafkaConfig, channel chan uuid.UUID) (Producer, Consumer) {

	config := createConfig()
	if err := createTopic(cfg, config); err != nil {
		utils.Logger.Fatal(err.Error())
	}
	producer := newProducer(cfg, config)
	consumer := newConsumer(cfg, config, channel)
	return producer, consumer

}
