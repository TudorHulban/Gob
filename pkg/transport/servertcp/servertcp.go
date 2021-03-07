package servertcp

import (
	"bufio"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	"main/pkg/processor"
	"main/pkg/seri/serigob"
	"main/pkg/transport/configtransport"

	"github.com/TudorHulban/log"

	"github.com/pkg/errors"
)

// Server Structure consolidating methods for a TCP server.
type Server struct {
	configtransport.Cfg
	serverStopping bool
	comms          chan []byte
	stop           chan struct{} // stop channel

	P *processor.Proc
	L *log.LogInfo
}

// NewServer Constructor creating a server with default or given configuration.
func NewServer(cfg *configtransport.Cfg) (*Server, error) {
	var g serigob.MGob
	logger := log.New(log.DEBUG, os.Stdout, true)

	p, err := processor.NewProc(processor.GobProcessing(g), processor.Logger(logger))
	if err != nil {
		return nil, errors.WithMessage(err, "could not create default transport configuration")
	}

	result := &Server{
		comms: make(chan []byte),
		stop:  make(chan struct{}),
		P:     p,
		L:     logger,
	}

	if cfg == nil {
		result.Cfg = *configtransport.NewDefaultConfiguration()
		return result, nil
	}

	result.Cfg = *cfg
	return result, nil
}

// Serve Method to be invoked for serving connections.
func (s *Server) Serve() (<-chan []byte, chan struct{}, error) {
	go s.listen()
	go s.listenStop()
	return s.comms, s.stop, nil
}

// listen Private method
func (s *Server) listen() error {
	s.L.Printf("listening on IP:%s, port:%d", s.IP, s.Port)

	listener, errListen := net.Listen(s.Protocol, s.Cfg.Socket())
	if errListen != nil {
		return errors.WithMessage(errListen, "could not start TCP server")
	}

	for {
		conn, err := listener.Accept() // blocks until new connection.
		if err != nil {
			s.L.Info("failed to accept connection:", err)
			continue
		}

		if s.serverStopping {
			conn.Write([]byte("stopping service..."))
			conn.Close()
			continue
		}

		go s.handleConn(conn, s.comms)
	}

	s.L.Print("stopping listening...")

	return nil
}

// handleConn Private method for handling connection.
// Closes after client message.
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

		if len(message) > 0 {
			s.L.Debug("received: ", string(message))

			// invoking message processor
			s.P.MsgDecode(message)

			comms <- message
			break // closing connection after message. should we leave it open?
		}
	}

	conn.Write([]byte("thank you"))
}

// listenStop Used for stopping the server in configured termination time.
func (s *Server) listenStop() {
	sigInterupt := make(chan os.Signal)
	signal.Notify(sigInterupt, os.Interrupt)

	<-sigInterupt

	s.serverStopping = true

	s.L.Printf("stopping in %d seconds", s.TerminationSecs)
	time.Sleep(time.Duration(s.TerminationSecs) * time.Second)

	s.stop <- struct{}{}
}
