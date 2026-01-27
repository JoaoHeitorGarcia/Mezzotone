package services

import (
	"image"
	"log"
	"os"
)

func ConvertImageToString(filePath string) {
	imageFile, _ := os.Open(filePath)
	defer imageFile.Close()

	decodedImage, _, _ := image.Decode(imageFile)
	log.Print(decodedImage.Bounds())
}
