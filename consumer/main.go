package main

//requires: sudo apt install libasound2-dev pkg-config

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/streadway/amqp"
)

func main() {
	f, err := os.Open("woo.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	b := beep.NewBuffer(format)
	b.Append(streamer)

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"simple_queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			Play(b.Streamer(0, b.Len()))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-done
}

func Play(s beep.Streamer) {
	speaker.Play(s)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
