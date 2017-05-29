package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"

	"goshot/utility"

	"errors"

	"github.com/charles-d-burton/gphoto2go"
	"github.com/gin-gonic/gin"
)

func main() {
	mdnsServer := utility.BroadcastServer()
	defer mdnsServer.Shutdown()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/shot", func(c *gin.Context) {
		camera := new(gphoto2go.Camera)
		err := camera.Init()
		defer camera.Exit()
		if err > 0 {
			log.Println(gphoto2go.CameraResultToString(err))
			c.JSON(500, gin.H{
				"err": gphoto2go.CameraResultToString(err),
			})
		}
		//camera.Interrupt()
		cameraFilePath, err := camera.TriggerCaptureToFile()
		if err == 0 {

			cameraFileReader := camera.FileReader(cameraFilePath.Folder, cameraFilePath.Name)
			defer cameraFileReader.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(cameraFileReader)

			//camera.DeleteFile(cameraFilePath.Folder, cameraFilePath.Name)
			encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())
			c.JSON(200, gin.H{
				"image": encodedImage,
			})
			buf.Reset()
			//c.Data(200, "image/jpeg", buf.Bytes())

		} else {
			log.Println(gphoto2go.CameraResultToString(err))
			c.Error(errors.New(gphoto2go.CameraResultToString(err)))
			c.JSON(http.StatusInternalServerError, gphoto2go.CameraResultToString(err))
			return
		}

	})
	r.Run()
}
