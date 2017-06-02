// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"goshot/utility"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/charles-d-burton/gphoto2go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	serverPort      int
	serverInterface string
	quality         int
	mutex           = &sync.Mutex{}
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start serving the camera",
	Long: `Serve the application on the given port
or default to 8080`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Args: ", strings.Join(args, ""))
		serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "port on which the server will listen")
	serveCmd.Flags().StringVarP(&serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
	serveCmd.Flags().IntVarP(&quality, "quality", "q", 100, "quality of the returned JPEG image")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve() {
	mdnsServer := utility.BroadcastServer(serverPort)
	defer mdnsServer.Shutdown()

	r := gin.Default()

	//Capture an image and return the base64 encoded value
	r.GET("/shot", getShotJSON)
	r.GET("/rawshot", getRawShot)
	r.Run(":" + strconv.Itoa(serverPort))

}

//Return a Base64 encoded JSON object containing the picture
func getShotJSON(c *gin.Context) {

	data, err := snapPicture()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	if http.DetectContentType(data) == "image/jpeg" { //Check that we have some actual image data
		imgData := data
		if quality < 100 {
			imgData, _ = convertJPEG(data)
		}
		encodedImage := base64.StdEncoding.EncodeToString(imgData)
		c.JSON(200, gin.H{
			"image": encodedImage,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unkown File type",
		})
	}
}

//Return a raw JPEG instead of Base64 encoded JSON doc
func getRawShot(c *gin.Context) {
	data, err := snapPicture()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	if http.DetectContentType(data) == "image/jpeg" { //Check that we have some actual image data
		c.Data(200, "image/jpeg", data)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unkown File type",
		})
	}
}

/*
 * Interface with the camera and have it take a picture
 * return that photo as a byte array.
 */
func snapPicture() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()
	camera := new(gphoto2go.Camera)
	err := camera.Init()
	defer camera.Exit() //Make sure to exit the camera at the end

	cameraFilePath, err := camera.TriggerCaptureToFile()
	if err == 0 {
		defer camera.DeleteFile(cameraFilePath.Folder, cameraFilePath.Name)
		cameraFileReader := camera.FileReader(cameraFilePath.Folder, cameraFilePath.Name)
		defer cameraFileReader.Close()
		buf := new(bytes.Buffer)
		defer buf.Reset()
		_, err := buf.ReadFrom(cameraFileReader)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	log.Println(gphoto2go.CameraResultToString(err))
	return nil, errors.New(gphoto2go.CameraResultToString(err))
}

/*
 * Convert image quality of JPEG
 */

func convertJPEG(imgByte []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var opt jpeg.Options
	opt.Quality = 90
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = jpeg.Encode(writer, img, &opt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b.Bytes(), nil

}
