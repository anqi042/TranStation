package Gin

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func RunGin(channel chan []byte) {



	r := gin.Default()

	r.POST("/write", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "FUCKING YEAH",
		})
		b, err := ioutil.ReadAll(c.Request.Body)
		if err !=nil{

		}else{
			channel <- b
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
