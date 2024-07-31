package main

import (
	"task4/router"
)

func main() {

	r := router.SetupRouter()
	r.Run(":8080")
}
