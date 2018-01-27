package Middleware

import (
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/YiftachE/DockerHPC/Core"
)

type loggingMiddleware struct {
	file   *os.File
	logger log.Logger
}

func NewLoggingMiddleware(path string) (loggingMiddleware, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return loggingMiddleware{}, err
	}
	logger := log.Logger{Level: log.InfoLevel, Formatter: &log.JSONFormatter{}, Out: f}
	return loggingMiddleware{file: f, logger: logger}, nil
}
func (middleware *loggingMiddleware) Dispose() {
	middleware.file.Close()
}
func (middleware *loggingMiddleware) HandleMessage(packet Core.DataPacket) Core.DataPacket {
	data, err := packet.Data.GetData()
	if err != nil {
		middleware.logger.Fatal(err)
	} else {
		middleware.logger.WithFields(log.Fields{"queueName": packet.QueueName, "size": data, "identifier": packet.Identifier}).Info("Recieved packet")
	}
	return packet
}
