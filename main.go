package main

import "books/app"

func main() {
	app := &app.App{}
	app.Init()
	app.Start()
	app.GracefulShutdown()
}
