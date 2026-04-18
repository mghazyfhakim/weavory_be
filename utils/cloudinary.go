package utils

import (
	"context"

	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() {
	var err error

	cld, err = cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	if err != nil {
		panic(err)
	}
}

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   "weavory",
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}