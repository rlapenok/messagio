package kafka

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/utils"
	"gopkg.in/Shopify/sarama.v1"
)

type Producer interface {
	Close() error
	readErrors()
	readSuccesses()
	WriteMessage(id uuid.UUID)
}

type prod struct {
	topic string
	wg    *sync.WaitGroup
	prod  sarama.AsyncProducer
}

func newProducer(cfg config.KafkaConfig, saramaConfig *sarama.Config) Producer {

	asyncProducer, err := sarama.NewAsyncProducer(cfg.Url, saramaConfig)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	prod := &prod{
		prod:  asyncProducer,
		topic: cfg.Topic,
		wg:    &sync.WaitGroup{},
	}
	prod.readErrors()
	prod.readSuccesses()
	return prod
}

func (p prod) Close() error {
	err := p.prod.Close()
	p.wg.Wait()
	return err
}

func (p prod) readErrors() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		utils.Logger.Info("Start reading Producer Errors channel")
		for err := range p.prod.Errors() {
			kafkaMessage := err.Msg
			logMsg := fmt.Sprintf("Producer: message was sent unsuccessfully to topic:%s, partition:%d with offset:%d with Error: %v", kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset, err)
			utils.Logger.Error(logMsg)
		}
	}()
}
func (p prod) readSuccesses() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		utils.Logger.Info("Start reading Producer Successes channel")
		for ok := range p.prod.Successes() {
			logMsg := fmt.Sprintf("Producer: message was sent successfully to topic:%s, partition:%d with offset:%d ", ok.Topic, ok.Partition, ok.Offset)
			utils.Logger.Info(logMsg)
		}
	}()
}

func (p prod) WriteMessage(id uuid.UUID) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		msg, err := id.MarshalBinary()
		if err != nil {
			logMsg := fmt.Sprintf("Cannot convert request from client with id:%v for message for kafka", id)
			utils.Logger.Error(logMsg)
		} else {
			p.prod.Input() <- &sarama.ProducerMessage{
				Topic: p.topic,
				Value: sarama.ByteEncoder(msg),
			}
			utils.Logger.Info("Producer: message was send")

		}
	}()
}
