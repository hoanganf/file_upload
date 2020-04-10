package main

import (
	"github.com/hoanganf/file_upload/src"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	bean, err := src.InitBean()
	if err != nil {
		log.Fatalln("can not create bean", err)
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/paste", bean.FileUploadService.Get)

	router.Static("/paste/view", bean.FileUploadService.Directory)
	router.POST("/paste", bean.FileUploadService.Post)
	router.Run(":8083")
}
