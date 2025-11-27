package consumer

import (
	"context"
	"encoding/json"
	"log-service/env"
	"log-service/src/client"
	"log-service/src/common/log"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Kafka struct {
	kafka *kgo.Client
	cfg   *env.Consumer
	ctx   context.Context
	es    *client.Elasticsearch
}

func NewKafka(ctx context.Context, cfg *env.Consumer) *Kafka {
	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBrokers...),
		kgo.ConsumerGroup(cfg.KafkaGroupID),
		kgo.ConsumeTopics(cfg.KafkaTopic),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	adm := kadm.NewClient(kafka)
	ensureTopicExists(adm, cfg.KafkaTopic)
	ensureTopicExists(adm, cfg.KafkaDLQTopic)

	es := client.NewElasticsearch(cfg.ElasticsearchURLs, cfg.ElasticsearchIndex)

	return &Kafka{
		ctx:   ctx,
		cfg:   cfg,
		kafka: kafka,
		es:    es,
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

func (k *Kafka) Run() {
	log.Logger.Info("starting kafka consumer")

	for {
		select {
		case <-k.ctx.Done():
			log.Logger.Info("stopping kafka consumer")
			return
		default:
			fetches := k.kafka.PollFetches(k.ctx)

			fetches.EachError(func(s string, i int32, err error) {
				log.Logger.Error("error fetching", zap.String("topic", s), zap.Int32("partition", i), zap.Error(err))
			})

			fetches.EachRecord(func(r *kgo.Record) {
				if err := k.es.Index(r.Value); err != nil {
					k.sendToDLQ(r, err)
					return
				}
			})

			if err := k.kafka.CommitUncommittedOffsets(context.Background()); err != nil {
				log.Logger.Error("error committing offsets", zap.Error(err))
			}
		}
	}
}

func (k *Kafka) sendToDLQ(r *kgo.Record, err error) {
	jsonData, _ := json.Marshal(map[string]any{
		"topic":     r.Topic,
		"partition": r.Partition,
		"offset":    r.Offset,
		"value":     string(r.Value),
		"error":     err.Error(),
		"timestamp": time.Now().Format(time.RFC3339),
	})

	k.kafka.Produce(context.Background(), &kgo.Record{
		Topic: k.cfg.KafkaDLQTopic,
		Value: jsonData,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Logger.Error("error sending to DLQ", zap.Error(err))
		}
	})

	log.Logger.Warn("sent to DLQ", zap.String("topic", k.cfg.KafkaDLQTopic))
}

func (k *Kafka) Close() {
	k.es.Close()
	k.kafka.Close()
}
