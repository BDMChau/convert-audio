package shared

import (
	"audio-convert/constants"
)

func CheckFileTypeIsNotFlacByEx(fileExtension string) bool {

	for _, value := range constants.FlacFileExtensions {
		if fileExtension == value {
			return true
		}
	}

	return false
}
