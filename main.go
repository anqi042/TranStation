package main

import (
	"./Gin"
	"./Parser"


)




func main()  {
	Parser.Init_logger()

	channelIn := make(chan  []byte,Parser.CHANNEL_HTTP_IN_LEN)
	Parser.RaiseParser(channelIn)
	Gin.RunGin(channelIn)



}

