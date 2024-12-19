package repositories

import (
	"log"
	"os"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/joho/godotenv"
)

var ik *imagekit.ImageKit

func InitImageKit() *imagekit.ImageKit {
	if ik != nil {
		return ik
	}

	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ik = imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		PublicKey:   os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		UrlEndpoint: os.Getenv("IMAGEKIT_URL_ENDPOINT"),
	})

	return ik
}
