package consumer

import (
	"context"
	"order-service/env"
	"order-service/src/broker/service"
	"order-service/src/common/log"
	"order-service/src/factory"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn                 *amqp091.Connection
	midtransTxNotifChann *amqp091.Channel
	ctx                  context.Context
	service              *service.RabbitMQ
}

func NewRabbitMQ(ctx context.Context, f *factory.Factory) *RabbitMQ {
	conn := connect()
	midtransTxNotifChann := setupChannel(conn, "notification", "midtrans-tx", "midtrans-tx")

	return &RabbitMQ{
		conn:                 conn,
		midtransTxNotifChann: midtransTxNotifChann,
		ctx:                  ctx,
		service:              service.NewRabbitMQ(f),
	}
}

func (c *RabbitMQ) checkCloseConnection(closeChann chan *amqp091.Error) {
	for err := range closeChann {
		if err != nil {
			c.reconnect()
		}
	}
}

func (c *RabbitMQ) reconnect() {
	if c.conn != nil && !c.conn.IsClosed() {
		return
	}

	c.conn = connect()
	c.midtransTxNotifChann = setupChannel(c.conn, "notification", "midtrans-tx", "midtrans-tx")

	go c.MidtransTx()

	log.Logger.Info("reconnect consumer success")
}

func (c *RabbitMQ) Close() {
	if err := c.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}

func setupChannel(conn *amqp091.Connection, exchangeName, queueName, keyName string) *amqp091.Channel {
	chann, err := conn.Channel()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	if err := chann.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	queue, err := chann.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	if err := chann.QueueBind(queue.Name, keyName, exchangeName, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	return chann
}

func connect() *amqp091.Connection {
	for {
		conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
			Heartbeat: 10 * time.Second,
		})

		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		return conn
	}

}
