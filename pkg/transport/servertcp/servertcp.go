package servertcp

import (
	"bufio"
	"io"
	"net"

	"main/pkg/transport/configtransport"

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
			s.L.Info("failed to accept connection:", err)
			continue
		}
		go s.handleConn(conn, s.comms)
	}

	return nil
}

// handleConn Connection closes after client message.
func (s *Server) handleConn(conn net.Conn, comms chan []byte) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, errReader := reader.ReadBytes(byte('\n'))
		if errReader != nil {
			if errReader != io.EOF {
				s.L.Info("failed to read data", errReader)
			}
			return
		}

		// closing connection after message. should we leave it open?
		if len(message) > 0 {
			s.L.Debug("received: ", string(message))
			comms <- message
			break
		}
	}

	conn.Write([]byte("thank you"))
}
