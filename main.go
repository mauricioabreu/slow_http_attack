package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	var conn net.Conn
	var err error
	var timerChan <-chan time.Time
	timer := time.NewTimer(10 * time.Second)
	timerChan = timer.C

loop:
	for {
		if conn != nil {
			conn.Close()
		}

		conn, err = net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		fmt.Fprint(conn, "POST / HTTP/1.1\r\n")
		for {
			select {
			case <-time.After(1 * time.Second):
				log.Info().Msg("sending body...")
				if _, err := fmt.Fprintf(conn, "a"); err != nil {
					log.Error().Err(err).Msg("")
					continue loop
				}
			case <-timerChan:
				timer.Reset(10 * time.Second)
				fmt.Fprintf(conn, "\r\n\r\n")
				io.ReadAll(conn)
				log.Info().Msg("closing connection...")
				conn.Close()
				continue loop
			}
		}
	}
}
