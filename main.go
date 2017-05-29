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
	//Broadcast where the service is
	mdnsServer := utility.BroadcastServer()
	defer mdnsServer.Shutdown()

	r := gin.Default()

	//Capture an image and return the base64 encoded value
	r.GET("/shot", func(c *gin.Context) {
		camera := new(gphoto2go.Camera)
		err := camera.Init()
		defer camera.Exit() //Make sure to exit the camera at the end
		if err > 0 {
			log.Println(gphoto2go.CameraResultToString(err))
			c.JSON(500, gin.H{
				"error": gphoto2go.CameraResultToString(err),
			})
		}
		//camera.Interrupt()
		cameraFilePath, err := camera.TriggerCaptureToFile()
		if err == 0 {

			cameraFileReader := camera.FileReader(cameraFilePath.Folder, cameraFilePath.Name)
			defer cameraFileReader.Close()
			buf := new(bytes.Buffer)
			defer buf.Reset()
			_, err := buf.ReadFrom(cameraFileReader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			if http.DetectContentType(buf.Bytes()) == "image/jpeg" { //Check that we have some actual image data
				encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())
				c.JSON(200, gin.H{
					"image": encodedImage,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Unkown File type",
				})
			}
			camera.DeleteFile(cameraFilePath.Folder, cameraFilePath.Name)

			//c.Data(200, "image/jpeg", buf.Bytes())

		} else {
			log.Println(gphoto2go.CameraResultToString(err))
			c.Error(errors.New(gphoto2go.CameraResultToString(err)))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gphoto2go.CameraResultToString(err),
			})
			return
		}

	})
	r.Run()
}
