package main

import "github.com/antoniokichaev/hezzl-collector/internal/app"

const configPath = "./config/config.yaml"

func main() {
	app.Run(configPath)
}
