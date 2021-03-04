package servertcp

type Cfg struct {
	Port     uint
	IP       string
	Protocol string
}

type Server struct {
	Cfg
	comms chan []byte
	stop  chan struct{}
}

func defaultConfiguration() *Cfg {
	return &Cfg{
		Port:     8080,
		IP:       "localhost",
		Protocol: "tcp",
	}
}

func NewServer(cfg *Cfg) (*Server, error) {
	result := new(Server)
	result.comms = make(chan []byte)
	result.stop = make(chan struct{})

	if cfg == nil {
		result.Cfg = *defaultConfiguration()
		return result, nil
	}

	result.Cfg = *cfg
	return result, nil
}

func (s *Server) Serve() (<-chan []byte, chan<- struct{}, error) {
	return s.comms, s.stop, nil
}
