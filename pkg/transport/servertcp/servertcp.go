package servertcp

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

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

func (s *Server) Serve() (<-chan []byte, chan struct{}, error) {
	return s.comms, s.stop, nil
}

func (s *Server) listen() error {
	listener, errListen := net.Listen(s.Protocol, s.IP+":"+strconv.FormatUint(uint64(s.Port), 10))
	if errListen != nil {
		return errors.WithMessage(errListen, "could not start TCP server")
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("failed to accept connection:", err)
				continue
			}
			go handleConn(conn, s.comms)
		}
	}()

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
