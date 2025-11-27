package consumer

import (
	"context"
	"email-service/env"
	"email-service/src/broker/publisher"
	"email-service/src/broker/service"
	"email-service/src/common/dto/request"
	"email-service/src/common/log"
	"email-service/src/factory"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn                 *amqp091.Connection
	otpChann             *amqp091.Channel
	ctx                  context.Context
	service              *service.RabbitMQ
	rabbitMQLogPublisher *publisher.Kafka
}

func NewRabbitMQ(ctx context.Context, f *factory.Factory) *RabbitMQ {
	conn := connect()
	otpChann := setupChannel(conn, "user", "otp", "otp")

	return &RabbitMQ{
		conn:                 conn,
		otpChann:             otpChann,
		ctx:                  ctx,
		service:              service.NewRabbitMQ(f),
		rabbitMQLogPublisher: f.RabbitMQLogPublisher,
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
	c.otpChann = setupChannel(c.conn, "user", "otp", "otp")

	go c.Otp()

	log.Logger.Info("reconnect consumer success")
}

func (c *RabbitMQ) Close() {
	if err := c.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}

func (c *RabbitMQ) publishLog(body any, msg amqp091.Delivery, start time.Time, tag, queue string, err error) {
	req := request.RabbitMQLog{
		Timestamp:    time.Now().Format(time.RFC3339),
		Exchange:     msg.Exchange,
		RoutingKey:   msg.RoutingKey,
		Queue:        queue,
		ConsumerTag:  tag,
		Payload:      body,
		ContentType:  msg.ContentType,
		DeliveryMode: msg.DeliveryMode,
		AppID:        msg.AppId,
		Acked:        err == nil,
		Latency:      float64(time.Since(start).Nanoseconds()) / 1e6,
	}

	if reqID, ok := msg.Headers["request_id"].(string); ok {
		req.RequestID = reqID
	}

	if userID, ok := msg.Headers["user_id"].(string); ok {
		req.UserID = &userID
	}

	if err != nil {
		errStr := err.Error()
		req.Error = &errStr
	}

	if err := c.rabbitMQLogPublisher.Publish(req); err != nil {
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

	dlqExchangeName := exchangeName + "-dlx"
	if err := chann.ExchangeDeclare(dlqExchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	dlqQueueName := queueName + "-dlq"
	if _, err := chann.QueueDeclare(dlqQueueName, true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	if err := chann.QueueBind(dlqQueueName, keyName+"-dlq", dlqExchangeName, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	args := amqp091.Table{
		"x-dead-letter-exchange":    dlqExchangeName,
		"x-dead-letter-routing-key": keyName + "-dlq",
	}

	queue, err := chann.QueueDeclare(queueName, true, false, false, false, args)
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
