# Janky Flat IoT Humble Beginning

This is a project that allows my flatmates to click a button in their browser and have a weird noise play elsewhere in 
the flat.

It composes of three parts, that can all be run on separate machines:

* RabbitMQ
  * A popular open source message broker
* Producer
  * An http service written in Go that runs a little web UI
  * The web UI has a button that POSTs to the `/trigger` endpoint of the service
  * That endpoint pushes a message to a queue on RabbitMQ
* Consumer
  * A little Go program that listens to the aforementioned RabbitMQ queue
  * Whenever a message is consumed, it plays a little gnome sound

## Notes for running

### Producer

* Requires environment variable `RABBITMQ_ADDRESS` to be set to the address of the RabbitMQ service, like such: 
  * `amqp://username:password@server:5672/`
* Service runs on port 8080

Build docker image with:

```sh
docker build ./producer -t producer
```

Run with:

```sh
docker run -d -e RABBITMQ_ADDRESS=amqp://username:password@server:5672/ -p "8080:8080" producer
```

### Consumer

* To have the docker container play sound this is needed: `--device /dev/snd` 
* Requires `RABBITMQ_ADDRESS` like in the producer

```sh
docker build ./consumer -t consumer
```

Run with:

```sh
docker run -d -e RABBITMQ_ADDRESS=amqp://username:password@server:5672/ --device /dev/snd consumer
```

### RabbitMQ

You can use `rabbitmq/docker-compose.yml` to set up a quick n dirty RabbitMQ service

## Other notes

Project used to learn about RabbitMQ and writing simple services in golang.