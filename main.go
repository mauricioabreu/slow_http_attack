package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {

loop:
	for {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		var timerChan <-chan time.Time
		var timer *time.Timer
		timer = time.NewTimer(10 * time.Second)
		timerChan = timer.C

		fmt.Fprint(conn, "POST / HTTP/1.1\r\n")
		for {
			select {
			case <-time.After(1 * time.Second):
				if timer != nil {
					timer.Reset(10 * time.Second)
				}
				log.Info().Msg("sending body...")
				_, err := fmt.Fprintf(conn, "a")
				if err != nil {
					log.Error().Err(err).Msg("")
					continue loop
				}
			case <-timerChan:
				fmt.Fprintf(conn, "\r\n\r\n")
				ioutil.ReadAll(conn)
				log.Info().Msg("closing connection...")
				conn.Close()
			}
		}
	}
}
