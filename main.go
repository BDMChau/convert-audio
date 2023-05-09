package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var convertTo = ".mp3"

const maxGoroutines = 10

func handleConvert(file fs.FileInfo, guard <-chan struct{}) {
	pathFile := "./music/" + file.Name()

	splitName := strings.Split(file.Name(), ".")
	splitName[len(splitName)-1] =
		strings.Replace(splitName[len(splitName)-1], splitName[len(splitName)-1], convertTo, -1)

	newPathFile := "./output/" + strings.Join(splitName, "")

	err := exec.Command("ffmpeg", "-i", pathFile, "-ab", "320k", "-map_metadata", "0", "-id3v2_version", "3", newPathFile).Run()
	if err != nil {
		fmt.Println("Error:", file.Name(), err.Error())
	}

	<-guard
	fmt.Println("OK:", file.Name())

}

func handleReadFiles(f *os.File, wg *sync.WaitGroup, guard chan struct{}) {
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		guard <- struct{}{}
		wg.Add(1)
		go handleConvert(file, guard)

		defer wg.Done()
	}
	wg.Wait()

}

// //////////////
func main() {
	wg := new(sync.WaitGroup)
	guard := make(chan struct{}, maxGoroutines) // to limit Go-routines

	f, err := os.Open("./music")

	if err != nil {
		err := fmt.Errorf("Error: %q", err)
		fmt.Println(err)
		return
	}

	handleReadFiles(f, wg, guard)
}
