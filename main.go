package main

import (
	"./Gin"
	"./Parser"
	"os"
	"log"
)


func main()  {

	channelIn := make(chan  []byte,1000)
	Parser.RaiseParser(channelIn)
	Gin.RunGin(channelIn)



}

