package Middleware

import "github.com/YiftachE/DockerHPC/Core"

type Middleware interface {
	HandleMessage(packet Core.DataPacket) Core.DataPacket
	Dispose()
	}
