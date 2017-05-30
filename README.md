# Goshot

A REST wrapper for libgphoto2 in Golang

## Installation

```
go get github.com/charles-d-burton/goshot

```

## Requirements

You will also need libgphoto2 installed. If you are on Mac OS X, I recommend installing it with homebrew.
```
brew install libgphoto2

```
Linux
```
sudo apt install libgphoto2-dev

```
Also install the Go wrapper for libgphoto2 and gin
```
go get github.com/charles-d-burton/gphoto2go
go get github.com/gin-gonic/gin

```


## Usage

The goal of this program is to provide a generic REST interface to remotely control a supported 
camera from libgphoto2.  

### Setup

Connect a supported camera to a computer
Ensure that the computer does NOT mount the camera as a drive this will prevent the camera from initializing
In the goshot directory
```
go run main.go

```

### Taking a Photo
Make a curl or other GET request against localhost:8080/shot

The return JSON contains an "image" field that contains a Base64 encoded JPEG
