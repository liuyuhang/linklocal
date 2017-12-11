package rabbitmq

import (
	"time"
	"fmt"

	"github.com/ccwings/log"
	"github.com/streadway/amqp"
)

var (
	rb_host, rb_user, rb_passwd, rb_hostname string
	rabbitmq_jobs chan amqp.Delivery
	rabbitmqconn RabbitmqConn
	routingkey string
	queue string
)

const(
	exchange = "linklocal"
	modelname = "api"
)

type RabbitmqMsg struct {
	Type        string
	Message     string
	FromRouting string
}

type RabbitmqConn struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     bool
	done    chan error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Debug("%s: %s", msg, err)
	}
}

func GetChannel() (ch *amqp.Channel, err error) {
	if rabbitmqconn.tag{
		ch, err = rabbitmqconn.conn.Channel()
		return ch,nil
	}else{
		return nil, <-rabbitmqconn.done
	}
}

func InitConnetion(jobs chan amqp.Delivery)(err error){
	for{
		rabbitmqconn.conn,rabbitmqconn.channel,err = setup(jobs)
		if err !=nil{
			time.Sleep(5*time.Second)
		}else {
			rabbitmqconn.tag = true
			return
		}
	}
}

func InitAMQP(jobs chan amqp.Delivery, host, user, passwd, hostname string) (err error) {
	rb_host = host
	rb_user = user
	rb_passwd = passwd
	rb_hostname = hostname

	routingkey = modelname+"."+rb_hostname
	queue = modelname+"."+rb_hostname

	rabbitmq_jobs = jobs
	rabbitmqconn.tag = false
	InitConnetion(rabbitmq_jobs)
	go ConnTest(rabbitmq_jobs)
	go JobRouter(rabbitmq_jobs)
	return
}

func ConnTest(jobs chan amqp.Delivery)(err error){
	for{
		if rabbitmqconn.conn != nil{
			fmt.Println(<-rabbitmqconn.conn.NotifyClose(make(chan *amqp.Error)))
			rabbitmqconn.tag = false
		}
		rabbitmqconn.conn,rabbitmqconn.channel,err = setup(jobs)
		if err!= nil{
			log.Info("Reconnect RabbitMQ Faild")
			rabbitmqconn.tag = false
		}else{
			log.Info("Reconnect RabbitMQ Sucessed")
			rabbitmqconn.tag = true
		}
		time.Sleep(5*time.Second)
	}
}

func setup(jobs chan amqp.Delivery) (conn *amqp.Connection, ch *amqp.Channel, err error) {
	url := "amqp://" + rb_user + ":" + rb_passwd + "@" + rb_host + ":5672/"
	conn, err = amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	if err := ch.ExchangeDeclare(exchange, "topic",true,false,false,false,nil); err != nil {
		return nil, nil, err
	}
	if _,err := ch.QueueDeclare(routingkey, false,true,false,false,nil);err != nil{
		return nil, nil, err
	}
	if err := ch.QueueBind(queue, routingkey,exchange,false,nil); err != nil{
		return nil, nil, err
	}
	if err := ch.QueueBind(queue, modelname+"_cast.*", exchange, false, nil); err != nil{
		return nil, nil, err
	}
	deliveries, err := ch.Consume(queue, routingkey, true, false, false, false, nil)
	if err!= nil{
		return nil, nil, err
	}
	go DeliveryJobs(rabbitmq_jobs,deliveries)
	return conn,ch,nil
}



func DeliveryJobs(jobs chan amqp.Delivery,deliveries <-chan amqp.Delivery) (err error) {
	for d := range deliveries {
		jobs <- d
		log.Info("Got a msg")
	}
	return
}

