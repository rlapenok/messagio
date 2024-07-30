package kafka

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/utils"
	"gopkg.in/Shopify/sarama.v1"
)

var channelForClose chan struct{}

type Consumer interface {
	ConsumeMessage()
	Close() error
}

type consumer struct {
	wg          *sync.WaitGroup
	consumer    sarama.Consumer
	topic       string
	chanForRepo chan uuid.UUID
	count       uint
}

func newConsumer(cfg config.KafkaConfig, saramaConfig *sarama.Config, channel chan uuid.UUID) Consumer {
	cons, err := sarama.NewConsumer(cfg.Url, saramaConfig)
	if err != nil {
		utils.Logger.Fatal(errors.Join(ErrCreateConsumer, err).Error())
	}
	return &consumer{
		consumer:    cons,
		topic:       cfg.Topic,
		wg:          &sync.WaitGroup{},
		chanForRepo: channel}
}

func (c consumer) Close() error {
	channelForClose <- struct{}{}
	close(channelForClose)
	err := c.consumer.Close()
	c.wg.Wait()
	close(c.chanForRepo)
	return err

}

func (c consumer) ConsumeMessage() {
	//get partitions in topic
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		utils.Logger.Fatal(errors.Join(ErrGetPartitions, err).Error())
	}
	channelForClose = make(chan struct{}, len(partitions))
	c.count = uint(len(partitions))
	//range for partitions
	for _, partition := range partitions {
		c.wg.Add(1)
		go func() {
			wg := &sync.WaitGroup{}
			defer c.wg.Done()
			defer wg.Wait()
			partitionConsumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
			if err != nil {
				logMsg := fmt.Sprintf("Consumer: cannot create  PartitionConsumer for topic :%s, partition:%d with error:%v", c.topic, partition, err)
				utils.Logger.Error(logMsg)
			} else {
				wg.Add(3)
				go func() {
					logMsg := fmt.Sprintf("Consumer: start consume messages from partition:%d Messages channel", partition)
					utils.Logger.Info(logMsg)
					defer wg.Done()
					for msg := range partitionConsumer.Messages() {
						logMsg := fmt.Sprintf("Consumer: message was successfully received in topic:%s,partition:%d,offset:%d", msg.Topic, msg.Partition, msg.Offset)
						utils.Logger.Info(logMsg)
						data, err := bytesIntoUuid(msg.Value)
						if err != nil {
							logMsg := fmt.Sprintf("Connot convert to uuid, error:%v", err)
							utils.Logger.Error(logMsg)
						} else {
							c.chanForRepo <- data
						}

					}
				}()
				go func() {
					logMsg := fmt.Sprintf("Consumer: start consume messages from partition:%d Errors channel", partition)
					utils.Logger.Info(logMsg)
					defer wg.Done()
					for err := range partitionConsumer.Errors() {
						logMsg := fmt.Sprintf("Consumer: message was unsuccessfully received in topic:%s,partition:%d with error %v", err.Topic, err.Partition, err)
						utils.Logger.Error(logMsg)

					}
				}()
				go func() {

					utils.Logger.Info("Consumer: start consume Close channel")
					defer wg.Done()
					<-channelForClose
					utils.Logger.Info("Start shutdown partions consumer")
					if err := partitionConsumer.Close(); err != nil {
						utils.Logger.Error(err.Error())

					}
					utils.Logger.Info("PartitonConsumer shutdown success")
				}()
			}
		}()
	}
}

func bytesIntoUuid(data []byte) (uuid.UUID, error) {

	return uuid.FromBytes(data)
}
