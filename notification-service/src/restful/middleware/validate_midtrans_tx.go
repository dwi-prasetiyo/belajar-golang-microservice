package middleware

import (
	"crypto/sha512"
	"encoding/hex"
	"notification-service/env"
	"notification-service/src/common/dto/request"
	"notification-service/src/common/errors"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ValidateMidtransTx(c *fiber.Ctx) error {
	req := new(request.MidtransTx)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// SHA512(order_id + status_code + gross_amount + serverkey)
	key := req.OrderID + req.StatusCode + req.GrossAmount + env.Conf.Midtrans.ServerKey

	hash := sha512.New()
	hash.Write([]byte(key))
	sha512Hash := hash.Sum(nil)
	sha512HashString := hex.EncodeToString(sha512Hash)

	if sha512HashString != req.SignatureKey {
		return &errors.Response{
			HttpCode: 400,
			Message:  "invalid signature",
		}
	}

	req.OrderID = req.Metadata.OriginalOrderID

	c.Locals("midtransTx", req)

	return c.Next()
}
