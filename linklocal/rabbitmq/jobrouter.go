package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/ccwings/log"
	"github.com/streadway/amqp"
)

func JobRouter(jobs chan amqp.Delivery) (err error) {
	for {
		d := <-jobs
		log.Info("doing Job:", string(d.Body))
		if d.ReplyTo != "" {
			log.Debug(d.ReplyTo)
		}
		msg := string(d.Body)
		message_map := make(map[string]interface{}, 0)
		err = json.Unmarshal([]byte(msg), &message_map)
		if err != nil {
			log.Debug(err)
			continue
		}
		switch message_map["type"].(string) {
		case "monitor":
			switch message_map["function"].(string) {
			case "health_report":
				log.Debug("in health report")
				go func() {
					log.Debug("go health report")
				}()
			default:
				err = fmt.Errorf("Message Function Not Found.", message_map["function"].(string))
				continue
			}
		default:
			err = fmt.Errorf("Message Type:%s not found", message_map["type"].(string))
			continue
		}

	}

}
