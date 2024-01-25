package main

import "github.com/timohahaa/ewallet/internal/app"

const configFilePath = "./config/config.yaml"

func main() {
	app.Run(configFilePath)
}
