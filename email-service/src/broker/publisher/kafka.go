package publisher

import (
	"context"
	"encoding/json"
	"time"
	"email-service/env"
	"email-service/src/common/log"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Kafka struct {
	client *kgo.Client
	topic  string
}

func NewKafka(topic string) *Kafka {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(env.Conf.Kafka.Addr1, env.Conf.Kafka.Addr2, env.Conf.Kafka.Addr3),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.ProducerLinger(5*time.Second),
		kgo.RequiredAcks(kgo.LeaderAck()),
		kgo.DisableIdempotentWrite(),
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	adm := kadm.NewClient(client)
	ensureTopicExists(adm, topic)

	return &Kafka{
		client: client,
		topic:  topic,
	}
}

func ensureTopicExists(adm *kadm.Client, topic string) {
	resp, err := adm.CreateTopics(context.Background(), 3, 2, nil, topic)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for _, r := range resp {
		if r.Err != nil {
			if r.Err == kerr.TopicAlreadyExists {
				log.Logger.Info("Topic already exists", zap.String("topic", topic))
				continue
			}
			log.Logger.Fatal(r.Err.Error())
		}

		log.Logger.Info("Topic created", zap.String("topic", topic))
	}
}

func (k *Kafka) Publish(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	k.client.Produce(context.Background(), &kgo.Record{
		Topic: k.topic,
		Value: jsonData,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Logger.Error(err.Error())
		}
	})

	return nil
}

func (k *Kafka) Close() {
	if err := k.client.Flush(context.Background()); err != nil {
		log.Logger.Error(err.Error())
	}

	k.client.Close()
}
