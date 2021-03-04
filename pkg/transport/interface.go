package transport

type ITransportClient interface {
	Send(Cfg, payload []byte) error
}

type ITransportServer interface {
	Serve() (payload <-chan []byte, stopSignal chan<- struct{}, error)
}
