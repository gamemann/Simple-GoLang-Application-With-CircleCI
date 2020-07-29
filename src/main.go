package main

import (
	"./webserver"
)

func main() {
	webserver.StartServer("0.0.0.0", 4030)
}
