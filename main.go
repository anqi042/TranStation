package main

import (
	_"./Gin"
	_"./Parser"
)
func main()  {
	channelIn := make(chan  []byte,1000)
	Parser.RaiseParser(channelIn)
	Gin.RunGin(channelIn)



}

