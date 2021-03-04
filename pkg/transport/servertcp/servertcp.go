package servertcp

import (
	"bufio"
	"io"
	"log"
	"net"

	"main/pkg/transport/config"

	"github.com/pkg/errors"
)

type Server struct {
	configtransport.Cfg
	comms chan []byte
	stop  chan struct{}
}

func NewServer(cfg *configtransport.Cfg) (*Server, error) {
	result := new(Server)
	result.comms = make(chan []byte)
	result.stop = make(chan struct{})

	if cfg == nil {
		result.Cfg = *configtransport.NewDefaultConfiguration()
		return result, nil
	}

	result.Cfg = *cfg
	return result, nil
}

func (s *Server) Serve() (<-chan []byte, chan struct{}, error) {
	go s.listen()
	return s.comms, s.stop, nil
}

func (s *Server) listen() error {
	s.L.Infof("listening on IP:%s, port:%d", s.IP, s.Port)

	listener, errListen := net.Listen(s.Protocol, s.Cfg.Socket())
	if errListen != nil {
		return errors.WithMessage(errListen, "could not start TCP server")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept connection:", err)
			continue
		}
		go handleConn(conn, s.comms)
	}

	return nil
}

// handleConn Connection remains opened until client closes it.
// TODO: assert addition of idle timeout
func handleConn(conn net.Conn, comms chan []byte) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, errReader := reader.ReadBytes(byte('\n'))
		if errReader != nil {
			if errReader != io.EOF {
				log.Println("failed to read data", errReader)
			}
			return
		}

		comms <- message
		log.Println("received: ", string(message))
	}
}
