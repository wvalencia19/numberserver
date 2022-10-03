package internal

import (
	"bufio"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/netutil"

	"net"
	"numberserver/internal/config"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	host      string
	port      string
	maxConn   int
	validator Validator
}

func New(config config.Server, validator Validator) *Server {
	return &Server{
		host:      config.Host,
		port:      config.Port,
		maxConn:   config.MaxConn,
		validator: validator,
	}
}

func (srv *Server) Run(ctx context.Context, cancel context.CancelFunc, logChan chan<- int) {
	log.Info("starting server")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", srv.host, srv.port))
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("max connections %d", srv.maxConn)
	listener = netutil.LimitListener(listener, srv.maxConn)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Debug(err)
		}
	}(listener)

	go gracefullyShutdown(cancel, listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Debug(err)
			return
		}
		go func() {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					return
				}
			}(conn)

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					fmt.Printf("closing connection: %v\n", ctx.Err())
					stopServer(conn, listener)
					return
				default:
					input := scanner.Text()
					if srv.validator.TerminationInput(input) {
						cancel()
						stopServer(conn, listener)
						return
					}
					valid, val := srv.validator.ValidInput(input)
					if valid {
						logChan <- val
					} else {
						err := conn.Close()
						if err != nil {
							log.Debug(err)
							return
						}
					}
				}
				if err := scanner.Err(); err != nil {
					log.Printf("Error scanning %v.\n", err)
				}
			}
		}()
	}
}

func stopServer(conn net.Conn, listener net.Listener) {
	err := conn.Close()
	if err != nil {
		log.Debug(err)
	}
	err = listener.Close()
	if err != nil {
		log.Debug(err)
	}
}

func gracefullyShutdown(cancel context.CancelFunc, listener net.Listener) {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("gracefully shutdown, server exiting")
	defer cancel()
	err := listener.Close()
	if err != nil {
		log.Debug(err)
	}
}
