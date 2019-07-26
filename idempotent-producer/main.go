package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)


func main() {

	// Define Vars
	var (
		kafka_version = "2.12.2"
		maxRetry = 5
		brokerList = []string {"localhost:9092"}
		topic = "test"
	) 

	version, err := sarama.ParseKafkaVersion(kafka_version) // Parse kafka version required minimum 0.11.0
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}


	config := sarama.NewConfig()
	config.Version = version // specity kafka version to sarama producer

	config.Producer.Idempotent = true // enable idempotent producer mode
	config.Net.MaxOpenRequests = 1 // Idempotent producer allows only one open request
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll // idempotent producer works with acks=all only
	config.Producer.Retry.Max = maxRetry


	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	for i := 0; i<100; i++ {

		message := fmt.Sprintf("something cool ... %d", i)

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}

		partition, offset, err := producer.SendMessage(msg)
        	if err != nil {
                	panic(err)
        	}

		fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	}
}