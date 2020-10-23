package main

import (
	"goTestProj/API"
)

func main() {
	api.InitializeDB().SetListeners()
}