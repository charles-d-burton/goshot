package main

import (
	"bytes"
	"log"

	"github.com/charles-d-burton/gphoto2go"
	"github.com/gin-gonic/gin"
)

func main() {
	camera := new(gphoto2go.Camera)
	err := camera.Init()
	if err > 0 {
		log.Println(err)
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/shot", func(c *gin.Context) {
		camera.Interrupt()
		cameraFilePath, err := camera.TriggerCaptureToFile()
		if err == 0 {

			cameraFileReader := camera.FileReader(cameraFilePath.Folder, cameraFilePath.Name)
			defer cameraFileReader.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(cameraFileReader)
			camera.DeleteFile(cameraFilePath.Folder, cameraFilePath.Name)

			c.Data(200, "image/jpeg", buf.Bytes())

		}
		log.Println(gphoto2go.CameraResultToString(err))
	})
	r.Run()
}
