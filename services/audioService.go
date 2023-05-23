package services

import (
	"audio-convert/constants"
	"audio-convert/models"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const maxGoroutines = 10

var convertTo = ".mp3"
var wg = new(sync.WaitGroup)

const InputPath = "./music"
const outputPath = "./output"

func handleConvertFlacFile(folderName string, file fs.FileInfo, userInputs models.UserInputs, originalExtension string) (result string, err error) {
	var pathFile string = InputPath + "/" + file.Name()
	if len(folderName) > 0 {
		pathFile = InputPath + "/" + folderName + "/" + file.Name()
	}
	newFileName := strings.Replace(file.Name(), originalExtension, convertTo, -1)
	newPathFile := outputPath + "/" + newFileName

	err = exec.Command("ffmpeg", "-i", pathFile, "-ab", strconv.Itoa(userInputs.BitRate)+"k", "-map_metadata", "0", "-id3v2_version", "3", newPathFile).Run()
	// err = exec.Command("ffmpeg", "-i", pathFile, "-c:a", "alac", "-c:v", "copy", newPathFile).Run() // use this for  FLAC to ALAC

	if err != nil {
		return "", err
	}

	return "ok", nil
}

func handleConvertOtherFile(folderName string, file fs.FileInfo, userInputs models.UserInputs, originalExtension string) (result string, err error) {
	var pathFile string = InputPath + "/" + file.Name()
	if len(folderName) > 0 {
		pathFile = InputPath + "/" + folderName + "/" + file.Name()
	}

	_file, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}
	defer _file.Close()

	destination, err := os.OpenFile(outputPath+"/"+file.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if originalExtension == ".mp3" {
		_, err = io.Copy(destination, _file)
	}
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func handleConvert(folderName string, file fs.FileInfo, guard <-chan struct{}, userInputs models.UserInputs) {
	defer func() {
		wg.Done()
		<-guard
	}()

	originalExtension := filepath.Ext(file.Name())
	if !constants.CheckFileType01ByEx(originalExtension) {
		_, err := handleConvertOtherFile(folderName, file, userInputs, originalExtension)
		if err != nil {
			fmt.Println("Error:", file.Name(), err)
			return
		}
	} else {
		_, err := handleConvertFlacFile(folderName, file, userInputs, originalExtension)
		if err != nil {
			fmt.Println("Error:", file.Name(), err)
			return
		}
	}

	fmt.Println("OK:", file.Name())
	return
}

func HandleReadFiles(f *os.File, userInputs models.UserInputs) {
	guard := make(chan struct{}, maxGoroutines) // to limit Go-routines

	var files []fs.FileInfo

	folders, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, folder := range folders {
		if folder.IsDir() {
			_folder, err := os.Open(InputPath + "/" + folder.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer _folder.Close()

			_files, err := _folder.Readdir(0)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, file := range _files {
				guard <- struct{}{}

				wg.Add(1)
				go handleConvert(folder.Name(), file, guard, userInputs)
			}
		} else {
			files = append(files, folder) // now, folder is a file
		}
	}

	if len(files) > 0 {
		for _, file := range files {
			guard <- struct{}{}

			wg.Add(1)
			go handleConvert("", file, guard, userInputs)
		}
	}

	wg.Wait()
}
