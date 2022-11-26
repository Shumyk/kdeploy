//go:build exclude

package legacy

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

const (
	DIVIDER   = "     "
	SPACE     = " "
	SEPARATOR = "|"
)

func select_image_prompt() {
	outputFile := os.Args[1]
	rawImagesInfo := os.Args[2:]

	prompt := &survey.Select{
		Message: "select image to deploy",
		Options: formatImagesInfo(rawImagesInfo),
	}
	selectedImage := ""
	survey.AskOne(prompt, &selectedImage)

	selectedImage = strings.Replace(selectedImage, DIVIDER, SPACE, -1)
	os.WriteFile(outputFile, []byte(selectedImage), 0666)
}

func formatImagesInfo(imagesInfo []string) []string {
	var imagesInfoFormatted []string
	for _, info := range imagesInfo {
		imagesInfoFormatted = append(
			imagesInfoFormatted,
			strings.Replace(info, SEPARATOR, DIVIDER, -1),
		)
	}
	return imagesInfoFormatted
}
