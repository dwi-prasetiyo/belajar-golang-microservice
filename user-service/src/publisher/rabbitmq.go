package publisher

import (
	"encoding/json"
	"sync"
	"time"
	"user-service/env"
	"user-service/src/common/dto/request"
	"user-service/src/common/log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn       *amqp091.Connection
	chann      *amqp091.Channel
	closeChann chan *amqp091.Error
	mutex      sync.Mutex
}

func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
		Heartbeat: 10 * time.Second,
	})

	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	closeChann := conn.NotifyClose(make(chan *amqp091.Error, 1))
	chann := setupChannel(conn)

	return &RabbitMQ{
		conn:       conn,
		chann:      chann,
		closeChann: closeChann,
	}
}

func (r *RabbitMQ) reconnect() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.conn != nil && !r.conn.IsClosed() {
		return
	}

	for {
		conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		r.conn = conn
		r.closeChann = r.conn.NotifyClose(make(chan *amqp091.Error, 1))

		log.Logger.Info("RabbitMQ reconnected")
		break
	}

	r.chann = setupChannel(r.conn)
}

func (r *RabbitMQ) Publish(exchange, key string, data *request.RabbitMQMessage) error {
	if r.conn == nil || r.conn.IsClosed() {
		r.reconnect()
	}

	jsonData, err := json.Marshal(data.Message)
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		Headers: amqp091.Table{
			"request_id": data.RequestID,
			"user_id":    data.UserID,
		},
		AppId:       "user-service",
		ContentType: "application/json",
		Body:        jsonData,
	}

	if err := r.chann.Publish(exchange, key, true, false, msg); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) Close() {
	if err := r.chann.Close(); err != nil {
		log.Logger.Error(err.Error())
	}

	if err := r.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}

func setupChannel(conn *amqp091.Connection) *amqp091.Channel {
	chann, err := conn.Channel()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	if err := chann.ExchangeDeclare("user", "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	return chann
}
