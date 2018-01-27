package Core

type LazyPacketData struct {
	data    []byte
	didGet  bool
	getData func() ([]byte, error)
}

func NewPacketData(getData func() ([]byte, error)) *LazyPacketData {

	return &LazyPacketData{didGet: false, getData: getData}
}
func (packet *LazyPacketData) GetData() ([]byte, error) {
	if !packet.didGet {
		data, error := packet.getData()
		if error != nil {
			return nil, error
		}
		packet.data = data
		packet.didGet = true
	}
	return packet.data, nil
}

type DataPacket struct {
	Identifier string
	QueueName  string
	Data       *LazyPacketData
}
