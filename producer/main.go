package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type result struct {
	Testing int `JSON:"testing"`
}

func main() {
	router := gin.Default()
	router.StaticFile("/", "index.html")
	router.GET("/trigger", Trigger)

	router.Run("0.0.0.0:8080")
}

func Trigger(c *gin.Context) {
	sendMessageToQueue()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessageToQueue() {
	// Here we connect to RabbitMQ or send a message if there are any errors connecting.
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))

	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"simple_queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// We set the payload for the message.
	body := "Hey!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	// If there is an error publishing the message, a log will be displayed in the terminal.
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Congrats, sending message: %s", body)
}
