package middleware

import (
	"time"
	"user-service/src/common/dto/request"
	"user-service/src/common/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) Log(c *fiber.Ctx) error {

	now := time.Now()
	requestID := uuid.New().String()
	c.Locals("request_id", requestID)

	err := c.Next()

	statusCode := c.Response().StatusCode()
	var errMessage string
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			statusCode = e.Code
		} else if e, ok := err.(*errors.Response); ok {
			statusCode = e.HttpCode
		} else {
			statusCode = 500
		}

		errMessage = err.Error()
	}

	var userID *string
	if id, ok := c.Locals("user_id").(string); ok {
		userID = &id
	}

	m.restfulLogPublisher.Publish(request.RestfulLog{
		RequestID:  requestID,
		Timestamp:  time.Now().Format(time.RFC3339),
		Method:     c.Method(),
		Url:        c.OriginalURL(),
		Body:       c.Locals("request_body"),
		StatusCode: statusCode,
		ClientIP:   c.IP(),
		UserAgent:  c.Get("User-Agent"),
		UserID:     userID,
		Latency:    float64(time.Since(now).Nanoseconds()) / 1e6,
		Error:      errMessage,
	})

	return err
}
