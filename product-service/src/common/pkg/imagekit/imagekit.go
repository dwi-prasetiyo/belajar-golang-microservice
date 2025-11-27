package imagekit

import (
	"product-service/env"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/logger"
)

func New() *imagekit.ImageKit {
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  env.Conf.ImageKit.PrivateKey,
		PublicKey:   env.Conf.ImageKit.PublicKey,
		UrlEndpoint: env.Conf.ImageKit.UrlEndpoint,
	})

	ik.Logger.SetLevel(logger.ERROR)
	return ik
}
