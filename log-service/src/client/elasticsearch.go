package client

import (
	"context"
	"log-service/src/common/log"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

type Elasticsearch struct {
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
}

func NewElasticsearch(urls []string, index string) *Elasticsearch {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: urls,
	})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: client,
		Index:  index,
		OnError: func(ctx context.Context, e error) {
			log.Logger.Error(e.Error())
		},
	})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	return &Elasticsearch{
		client:      client,
		bulkIndexer: bulkIndexer,
	}
}

func (e *Elasticsearch) Index(data []byte) error {
	return e.bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action: "create",
		Body:   strings.NewReader(string(data)),
	})
}

func (e *Elasticsearch) Close() {
	if err := e.bulkIndexer.Close(context.Background()); err != nil {
		log.Logger.Error(err.Error())
	}
}
