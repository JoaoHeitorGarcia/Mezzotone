package services

import "os"

func ConvertImageToString() {
	filePath := getSharedVariable("selectedFile")
	file, err := os.ReadFile(filePath.(string))
	if err != nil {
		//TODO dont panic :D
	}
	//TODO ...
}
