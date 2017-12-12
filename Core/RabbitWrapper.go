package Core

import "github.com/streadway/amqp"

type RabbitConnection struct {
	conn         *amqp.Connection
	ErrorChannel chan error
}

func NewRabbitConnection(rabbitConnString string) (RabbitConnection, error) {
	conn, err := amqp.Dial(rabbitConnString)
	return RabbitConnection{conn, make(chan error, 1)}, err
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
	headers["name"] = packet.Identifier
	channel.Publish("", q.Name, false, false, amqp.Publishing{Headers: headers, Body: packet.Data})
}
