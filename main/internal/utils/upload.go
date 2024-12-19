package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"commerce/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func UploadImagesToImageKit(files []*multipart.FileHeader) ([]string, error) {
	var uploadResults []string
	ik := repositories.InitImageKit()

	for _, fileHeader := range files {
		if !isValidFileType(fileHeader.Filename) {
			return nil, fmt.Errorf("invalid file type: %s", fileHeader.Filename)
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileBytes := make([]byte, fileHeader.Size)
		_, err = file.Read(fileBytes)
		if err != nil {
			return nil, err
		}

		base64Image := base64.StdEncoding.EncodeToString(fileBytes)

		uploadParam := uploader.UploadParam{
			FileName: fileHeader.Filename,
		}

		ctx := context.Background()
		uploadResp, err := ik.Uploader.Upload(ctx, base64Image, uploadParam)
		if err != nil {
			return nil, err
		}

		uploadResults = append(uploadResults, uploadResp.Data.Url)
	}

	return uploadResults, nil
}

func isValidFileType(fileName string) bool {
	allowedTypes := []string{".jpg", ".jpeg", ".png"}
	for _, t := range allowedTypes {
		if strings.HasSuffix(strings.ToLower(fileName), t) {
			return true
		}
	}

	log.Printf("Unsupported file type: %s", fileName)
	return false
}

const maxUploadFiles = 10

func HandleFileUploads(c *fiber.Ctx) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	files := form.File["images"]
	if len(files) > maxUploadFiles {
		return nil, fmt.Errorf("you can upload a maximum of %d files", maxUploadFiles)
	}
	return UploadImagesToImageKit(files)
}
