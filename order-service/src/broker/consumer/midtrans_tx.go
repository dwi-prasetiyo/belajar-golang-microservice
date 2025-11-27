package consumer

import (
	"order-service/src/common/log"

	"github.com/rabbitmq/amqp091-go"
)

func (c *RabbitMQ) MidtransTx() {
	log.Logger.Info("starting midtrans tx consumer")

	closeChann := c.conn.NotifyClose(make(chan *amqp091.Error, 1))
	go c.checkCloseConnection(closeChann)

	midtransTxConsumer, err := c.midtransTxNotifChann.ConsumeWithContext(c.ctx, "midtrans-tx", "midtrans-tx-consumer", false, false, false, false, nil)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for {
		select {
		case msg := <-midtransTxConsumer:
			if c.conn == nil || c.conn.IsClosed() {
				return
			}

			if err := c.service.HandleMidtransTxNotif(c.ctx, msg.Body); err != nil {
				if err := msg.Nack(false, true); err != nil {
					log.Logger.Error(err.Error())
					continue
				}

				log.Logger.Error(err.Error())
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Logger.Error(err.Error())
				continue
			}

			log.Logger.Info("successfully processes midtrans tx notification")

		case <-c.ctx.Done():
			return
		}
	}
}
