package main

import (
	"encoding/json"
	"fmt"
	"github.com/ricanontherun/rabbit-producer/utils"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type options struct {
	url         string
	exchange    string
	routingKey  string
	message     []byte
	contentType string
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("rabbit-producer: %s, %s", message, err.Error())
	}
}

func failWithFlagError(err error, flagSet *utils.RequiredFlagSet) {
	if err != nil {
		fmt.Println(err.Error())
		flagSet.Usage()
		os.Exit(2)
	}
}

func getStdIn() ([]byte, error) {
	// Make sure piped stdin is even enabled before doing this.
	return ioutil.ReadAll(os.Stdin)
}

func receivingPipedInput() bool {
	info, err := os.Stdin.Stat()
	failOnError(err, "Failed to stat Stdin")
	return info.Mode()&os.ModeCharDevice == 0 && info.Size() > 0
}

func parseFlags() *options {
	requiredFields := []string{}

	// We only to require the message flag is we aren't being piped data.
	pipedInput := receivingPipedInput()
	if !pipedInput {
		requiredFields = append(requiredFields, "message")
	}

	flagSet := utils.NewRequiredFlagSet(requiredFields)

	url := flagSet.String("remoteUrl", "localhost", "Connection URL for RabbitMQ (no port)")
	exchange := flagSet.String("exchange", "", "Exchange")
	routingKey := flagSet.String("routingKey", "", "Routing Key")
	message := flagSet.String("message", "", "Message to send")
	contentType := flagSet.String("contentType", "", "Content type to publish as")

	err := flagSet.Parse()
	if err != nil {
		failWithFlagError(err, flagSet)
	}

	options := &options{
		url:         *url,
		exchange:    *exchange,
		routingKey:  *routingKey,
		contentType: *contentType,
	}

	if pipedInput {
		input, err := getStdIn()
		failOnError(err, "Failed to read from stdin")
		options.message = []byte(strings.Replace(string(input), "\n", "", -1))
	} else {
		options.message = []byte(*message)
	}

	return options
}

func guessContentType(data []byte) string {
	if json.Valid(data) {
		return "application/json"
	}

	return "text/plain"
}

func main() {
	opts := parseFlags()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s", opts.url))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to create RabbitMQ channel")
	defer channel.Close()

	contentType := opts.contentType
	if len(contentType) == 0 {
		contentType = guessContentType(opts.message)
	}

	err = channel.Publish(opts.exchange, opts.routingKey, false, false, amqp.Publishing{
		ContentType: contentType,
		Body:        opts.message,
	})
	failOnError(err, "Failed to publish message to profiles.new")
}
