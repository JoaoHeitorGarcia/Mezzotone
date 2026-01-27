package services

import (
	"os"
)

func ConvertImageToString(filePath string) {
	imageFile, _ := os.Open(filePath)
	defer imageFile.Close()

	_ = Logger().Info("Sucessfully Loaded: " + filePath)
	
	//decodedImage, _, _ := image.Decode(imageFile)
	//_ = Logger().Info(decodedImage.Bounds().String())
}
