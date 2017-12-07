package main

// Main Function
// Been used to create api port
// init function used to read config and init DB & Controller

import (
	"github.com/ccwings/log"
	"gopkg.in/macaron.v1"
	"strconv"
	"linklocal/controller"
	"linklocal/model"
	"linklocal/auth"
	"linklocal/utils"
)

var (
	m          = macaron.Classic()
	listenPort = 3000
)

func init() {
	// Read Config File
	config, err := utils.ReadConf()
	if err != nil {
		log.Error("Config File Read Error")
	}

	// Define Listen Port
	listenPortStr, _ := config.GetString("default", "listen")
	listenPort, _ = strconv.Atoi(listenPortStr)

	// Init DB
	dataDriver, _ := config.GetString("default", "data_driver")
	dataSource, _ := config.GetString("default", "data_source")
	model.InitModelDB(dataDriver, dataSource)

	// Init Controller
	controller.InitRouter(m)

	// Init Token For Test Env
	auth.InitDevToken()

}

func main() {
	log.SetOutputLevel(0)
	log.Info("Service Starting...")
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: true,
	}))
	m.Run(listenPort)
}
