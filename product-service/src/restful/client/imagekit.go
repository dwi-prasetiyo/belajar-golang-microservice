package client

import (
	"context"
	"encoding/base64"
	"os"
	ik "product-service/src/common/pkg/imagekit"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type ImageKit struct {
	ik *imagekit.ImageKit
}

func NewImageKit() *ImageKit {
	ik := ik.New()

	return &ImageKit{
		ik: ik,
	}
}

func (c *ImageKit) UploadFile(ctx context.Context, path string, fileName string) (*uploader.UploadResult, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	base64Str := base64.StdEncoding.EncodeToString(fileData)
	file := "data:image/jpeg;base64," + base64Str

	useUniqueFileName := true

	res, err := c.ik.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName:          fileName,
		UseUniqueFileName: &useUniqueFileName,
	})

	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *ImageKit) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.ik.Media.DeleteFile(ctx, fileID)
	return err
}
