package Core

const INSTREAMNAME = "input_stream"

type InputStream interface {
	Start()
	HandleExceptions(chan error)
	Insert(chan DataPacket)
}
