package main

import (
	_"fmt"
	"zensey.go.archetype/logger"
)

func main() {
	logger.Log_info("Hello world !")
	logger.Log_info_f("Hello world %d !\n", 123)
	return
}