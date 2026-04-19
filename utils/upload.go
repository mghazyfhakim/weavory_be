package utils

import (
	"mime/multipart"
)

func UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	url, err := UploadToCloudinary(file, fileHeader.Filename)
	if err != nil {
		return "", err
	}

	return url, nil
}