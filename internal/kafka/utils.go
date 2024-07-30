package kafka

import (
	"errors"
	"fmt"
	"time"

	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/utils"
	"gopkg.in/Shopify/sarama.v1"
)

var ErrClusterAdim = errors.New("cannot create NewClusterAdmin")
var ErrCreateTopic = errors.New("cannot create topic")
var ErrCreateProducer = errors.New("cannot creaye AsyncProducer")
var ErrCreateConsumer = errors.New("cannot creaye Consumer")
var ErrGetPartitions = errors.New("cannot get partitions for Consumer")
var ErrCreatePartitionConsumer = errors.New("cannot get partitions for PartitionConsumer")
var ErrCloseClinet = errors.New("kafka: tried to use a client that was closed")

func createConfig() *sarama.Config {
	//create new saramaConfig for Producer And Comsuner
	config := sarama.NewConfig()

	config.Version = sarama.MaxVersion
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Flush.Frequency = 1 * time.Millisecond
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	return config
}

func createTopic(cfg config.KafkaConfig, saramaConfig *sarama.Config) error {
	//create newClusterAdmin for create topic
	admin, err := sarama.NewClusterAdmin(cfg.Url, saramaConfig)
	if err != nil {
		return errors.Join(err, ErrClusterAdim)
	}
	//create topic details
	topicDetails := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: int16(1),
	}
	//min handle errors
	switch err := admin.CreateTopic(cfg.Topic, &topicDetails, true); {
	case err != nil:
		{
			if errors.Is(err, sarama.ErrTopicAlreadyExists) {
				utils.Logger.Info(fmt.Sprintf("topic: %s already exist", cfg.Topic))
				return nil
			}
			return errors.Join(err, ErrCreateTopic)
		}
	default:
		{
			utils.Logger.Info(fmt.Sprintf("topic: %s was created", cfg.Topic))
			return nil
		}
	}
}
