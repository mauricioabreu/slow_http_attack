package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	targetAddress = "localhost:8080"
	sendDelay     = 1 * time.Second
	timerDuration = 10 * time.Second
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var conn net.Conn
	var err error
	var timerChan <-chan time.Time
	timer := time.NewTimer(timerDuration)
	timerChan = timer.C

outerLoop:
	for {
		if conn != nil {
			conn.Close()
		}

		conn, err = net.Dial("tcp", targetAddress)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		fmt.Fprint(conn, "POST / HTTP/1.1\r\n")

		for {
			select {
			case <-time.After(sendDelay):
				log.Info().Msg("sending body...")

				if _, err := fmt.Fprintf(conn, "a"); err != nil {
					log.Error().Err(err).Msg("")
					continue outerLoop
				}
			case <-timerChan:
				timer.Reset(timerDuration)
				fmt.Fprintf(conn, "\r\n\r\n")
				io.ReadAll(conn)
				log.Info().Msg("closing connection...")
				conn.Close()
				continue outerLoop
			}
		}
	}
}
