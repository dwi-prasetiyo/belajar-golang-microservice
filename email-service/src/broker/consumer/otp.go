package consumer

import (
	"email-service/src/common/dto/request"
	"email-service/src/common/log"
	"encoding/json"
	"time"
)

func (c *RabbitMQ) Otp() {
	log.Logger.Info("starting otp consumer")

	otpConsumer, err := c.otpChann.ConsumeWithContext(c.ctx, "otp", "otp-consumer", false, false, false, false, nil)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for {
		select {
		case msg := <-otpConsumer:
			now := time.Now()

			var logBody request.SendOtp
			body := new(request.SendOtp)

			if err := json.Unmarshal(msg.Body, &body); err != nil {
				c.publishLog(logBody, msg, now, "otp-consumer", "otp", err)
				continue
			}

			logBody = *body
			logBody.Otp = "*******"

			if c.conn == nil || c.conn.IsClosed() {
				c.publishLog(logBody, msg, now, "otp-consumer", "otp", nil)
				return
			}

			if err := c.service.SendOtp(msg.Body); err != nil {
				if err := msg.Nack(false, false); err != nil {
					c.publishLog(logBody, msg, now, "otp-consumer", "otp", err)
					continue
				}

				c.publishLog(logBody, msg, now, "otp-consumer", "otp", err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				c.publishLog(logBody, msg, now, "otp-consumer", "otp", err)
				continue
			}

			c.publishLog(logBody, msg, now, "otp-consumer", "otp", nil)

		case <-c.ctx.Done():
			return
		}
	}
}
