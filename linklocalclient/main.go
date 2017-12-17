package main

// Main Function
// Been used to create api port
// init function used to read config and init DB & Controller

import (
	"github.com/ccwings/log"
	"gopkg.in/macaron.v1"
	"linklocal/utils"
	"linklocal/rabbitmq"
	"github.com/streadway/amqp"
)

var (
	m          = macaron.Classic()
	jobs   = make(chan amqp.Delivery, 10)
)

func init() {
	log.Info("Init function started")
	// Read Config File
	config, err := utils.ReadConf()
	if err != nil {
		log.Error("Config File Read Error")
	}


	rabbitmq_enable, _ := config.GetBool("default", "rabbitmq_enable")
	if rabbitmq_enable {
		log.Info("Init AMQP...")
		rabbitmq_server, _ := config.GetString("default", "rabbitmq_server")
		uuid, _ := config.GetString("default", "uuid")
		err = rabbitmq.InitAMQP(jobs, rabbitmq_server, "linklocal", "zaq12wsxUu125478521", uuid)
		if err != nil {
			log.Debug(err)
		}
	}

	log.Info("Init function finished")

}

func main() {
	log.SetOutputLevel(0)
	log.Info("Service Starting...")
	forever := make(chan bool)
	log.Info("main Server Start")
	log.Info("Starting service...")
	<-forever
	log.Info("End service...")
}
