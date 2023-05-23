package constants

func CheckFileType01ByEx(fileExtension string) bool {
	fileType01 := []string{".flac", ".wav"}

	for _, value := range fileType01 {
		if fileExtension == value {
			return true
		}
	}

	return false
}
