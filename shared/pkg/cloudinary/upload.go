package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/api"
)

func Credentials() (*cloudinary.Cloudinary, context.Context) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")

	CLOUDINARY_URL := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)

	cld, _ := cloudinary.NewFromURL(CLOUDINARY_URL)
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadImage(cld *cloudinary.Cloudinary, ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       fmt.Sprintf("%s_%v", fileHeader.Filename, time.Now()),
		UniqueFilename: *api.Bool(false),
		Overwrite:      *api.Bool(true),
	})
	if err != nil {
		fmt.Println("error")
		return "", err
	}

	return resp.SecureURL, nil
}

func CloudinaryUploadImage(r *http.Request) (string, error) {
	// Retrieve the uploaded image file
	file, fileHeader, err := r.FormFile("image")

	if file == nil {
		return "image.jpeg", nil
	}

	if err != nil {
		return "", err
	}
	defer file.Close()

	// File Upload
	cld, ctx := Credentials()
	imageUrl, err := UploadImage(cld, ctx, file, fileHeader)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return imageUrl, nil
}
