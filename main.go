package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charles-d-burton/gphoto2go"
)

func main() {
	fmt.Printf("Hello world.\n")
	camera := new(gphoto2go.Camera)
	err := camera.Init()
	if err > 0 {
		log.Println(err)
	}
	cameraFilePath, err := camera.TriggerCaptureToFile()
	cameraFileReader := camera.FileReader(cameraFilePath.Folder, cameraFilePath.Name)
	os.Mkdir("/tmp/testing", os.FileMode(0777))
	fileWriter, err1 := os.Create("/tmp/testing/test.jpg")
	if err1 != nil {
		log.Println(err)
	}
	io.Copy(fileWriter, cameraFileReader)
	cameraFileReader.Close()
	/*folders := camera.RListFolders("/")
	for _, folder := range folders {
		files, _ := camera.ListFiles(folder)
		for _, fileName := range files {
			cameraFileReader := camera.FileReader(folder, fileName)
			os.Mkdir("/tmp/testing", os.FileMode(0777))
			fileWriter, err := os.Create("/tmp/testing/" + fileName)
			if err != nil {
				log.Println(err)
			}
			io.Copy(fileWriter, cameraFileReader)
			cameraFileReader.Close()
		}
	}*/
}
