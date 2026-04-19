package utils

import (
	"context"
	"log"
	"fmt"

	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() {
	var err error

	cloud := os.Getenv("CLOUDINARY_CLOUD_NAME")
	key := os.Getenv("CLOUDINARY_API_KEY")
	secret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloud == "" || key == "" || secret == "" {
		log.Fatal("❌ Cloudinary ENV NOT SET")
	}

	cld, err = cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		log.Fatal("❌ Failed init Cloudinary:", err)
	}

	log.Println("✅ Cloudinary initialized:", cloud)
}

func UploadToCloudinary(file multipart.File, filename string) (string, error) {

	if cld == nil {
		return "", fmt.Errorf("cloudinary not initialized")
	}

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