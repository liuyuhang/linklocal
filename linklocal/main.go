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
	"linklocal/rabbitmq"
	"github.com/streadway/amqp"
)

var (
	m          = macaron.Classic()
	listenPort = 3000
	jobs   = make(chan amqp.Delivery, 10)
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


	rabbitmq_enable, _ := config.GetBool("default", "rabbitmq_enable")
	if rabbitmq_enable {
		log.Info("Init AMQP...")
		rabbitmq_server, _ := config.GetString("default", "rabbitmq_server")
		uuid, _ := config.GetString("default", "uuid")
		err = rabbitmq.InitAMQP(jobs, rabbitmq_server, "linklocal", "zaq12wsx", uuid)
		if err != nil {
			log.Debug(err)
		}
	}

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
