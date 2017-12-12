package main

type Message struct{}
type Response struct{}
type DataPacket struct {
	Message
}

type Middleware interface {
	handleMessage(Message) Message
}

type DataPacketWorker interface {
	receiveMessage(DataPacket) Response
}

type LoggingMiddleware struct{}

func (LoggingMiddleware) handleMessage(Message) Message {
	panic("implement me")
}
