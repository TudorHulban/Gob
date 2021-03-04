package transport

type ITransportClient interface {
	Send(payload []byte) (string, error)
}

type ITransportServer interface {
	Serve() (<-chan []byte, chan<- struct{}, error) // communication and stop channels
}
