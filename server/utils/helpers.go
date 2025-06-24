package utils

import (
	"bytes"
	"errors"
	"image"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsImage(fileHeader *multipart.FileHeader) bool {
	fileType := fileHeader.Header.Get("Content-Type")
	return strings.HasPrefix(fileType, "image/")
}

func ExtractFileAndUpload(req *http.Request, fieldName string) (string, int, error) {
	file, fileHeader, err := req.FormFile(fieldName)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("File missing or invalid.")
	}

	if !IsImage(fileHeader) {
		return "", http.StatusBadRequest, errors.New("Only images are allowed.")
	}

	imgBytes, err := compressImage(file)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("Failed to decode the image.")
	}

	filePath, err := UploadFileToS3(imgBytes, fileHeader)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("Failed to upload the image.")
	}

	return filePath, 0, nil
}

func compressImage(file multipart.File) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.New("Failed to decode the image.")
	}

	img = imaging.Resize(img, 512, 0, imaging.Lanczos)

	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{
		Lossless: false,
		Quality:  75,
	})
	if err != nil {
		return nil, errors.New("Failed to decode the image.")
	}

	return buf.Bytes(), nil
}

func GetOtherParticipant(userID bson.ObjectID, participants [2]bson.ObjectID) bson.ObjectID {
	if userID.Hex() == participants[0].Hex() {
		return participants[1]
	}
	return participants[0]
}
