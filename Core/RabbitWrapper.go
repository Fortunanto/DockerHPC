package Core

import (
	"time"

	"github.com/YiftachE/DockerHPC/Core/Interfaces"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	conn         *amqp.Connection
	ErrorChannel chan error
	BackingStore Interfaces.BackingStore
}

func NewRabbitConnection(rabbitConnString string, store Interfaces.BackingStore) (RabbitConnection, error) {
	conn, err := amqp.Dial(rabbitConnString)
	return RabbitConnection{conn, make(chan error, 1), store}, err
}
func (rabbitConn RabbitConnection) Close() {
	rabbitConn.conn.Close()
}
func (rabbitConn RabbitConnection) SendMessage(queueName string, packet DataPacket) {
	channel, chanErr := rabbitConn.conn.Channel()
	if chanErr != nil {
		rabbitConn.ErrorChannel <- chanErr
		return
	}
	q, queueErr := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if queueErr != nil {
		rabbitConn.ErrorChannel <- queueErr
		return
	}
	headers := make(map[string]interface{})
	headers["identifier"] = packet.Identifier
	headers["receiverName"] = packet.QueueName
	data, errorAcquireData := packet.Data.GetData()
	if errorAcquireData != nil {
		rabbitConn.ErrorChannel <- errorAcquireData
		return
	}
	storageIdentifier, storeErr := rabbitConn.BackingStore.InsertValue(data , packet.Identifier, time.Now().String())
	if storeErr != nil {
		rabbitConn.ErrorChannel <- storeErr
		return
	}
	headers["storageIdentifier"] = storageIdentifier
	channel.Publish("", q.Name, false, false, amqp.Publishing{Headers: headers})
}

func (rabbitConn RabbitConnection) ReceiveMessage(queueName string) (chan DataPacket, error) {
	messages := make(chan DataPacket, 1)
	channel, chanErr := rabbitConn.conn.Channel()
	if chanErr != nil {
		return messages, chanErr
	}
	_, queueErr := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if queueErr != nil {
		return messages, queueErr
	}
	deliveries, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return messages, err
	}
	go func() {
		for msg := range deliveries {
			queueName := msg.Headers["receiverName"].(string)
			identifier := msg.Headers["identifier"].(string)
			redisIdentifier := msg.Headers["storageIdentifier"].(string)
			getData := func() ([]byte, error) {
				data, backingStoreError := rabbitConn.BackingStore.GetValue(redisIdentifier)
				return data.([]byte), backingStoreError
			}
			messages <- DataPacket{QueueName: queueName, Identifier: identifier, Data: NewPacketData(getData)}
		}
	}()
	return messages, nil
}
