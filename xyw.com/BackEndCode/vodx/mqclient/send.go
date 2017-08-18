package mqclient

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"xutils/xerr"

	"vodx/mqclient/mqview"
)

func Send(p *MQClient) {
	defer xerr.CatchPanic()

	ch, err := p.pConn.Channel()
	xerr.ThrowPanic(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	xerr.ThrowPanic(err)

	body := "hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body)})
	xerr.ThrowPanic(err)
}

// panic
func SendMessage(qName string, msg mqview.MQMessage, pClient *MQClient) {
	ch, err := pClient.pConn.Channel()
	xerr.ThrowPanic(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	xerr.ThrowPanic(err)

	body, _ := json.Marshal(&msg)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body})
	xerr.ThrowPanic(err)
}
