package mqclient

import (
	"log"

	"xutils/xerr"
)

func (p *MQClient) Receive() {
	defer xerr.CatchPanic()

	channel, err := p.pConn.Channel()
	xerr.ThrowPanic(err)
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	xerr.ThrowPanic(err)

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	xerr.ThrowPanic(err)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
