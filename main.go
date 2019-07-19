package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jaysoncena/bla/receiver"
	"github.com/jaysoncena/bla/sender"
)

var (
	flag         = kingpin.New("bla", "Multicast File Transfer")
	senderMode   = flag.Flag("sender", "Sender Mode").Bool()
	receiverMode = flag.Flag("receiver", "Receiver Mode").Bool()
	address      = flag.Flag("addr", "IPv4 Address (Listen IP or Target IP)").Required().String()
)

func init() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	kingpin.MustParse(flag.Parse(os.Args[1:]))

	if !*senderMode && !*receiverMode {
		log.Error().Msg("Please choose either sender or receiver mode")
	}
}

func main() {
	if *senderMode {
		s, err := sender.NewSender(*address)
		if err != nil {
			log.Panic().Msg(err.Error())
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Input: ")
			input, _, err := reader.ReadLine()
			if err != nil {
				log.Panic().Msg(err.Error())
			}

			b := make([]byte, 256)
			copy(b, input)
			_, err = s.Conn.Write(b)
			if err != nil {
				log.Panic().Msg(err.Error())
			}
		}

	} else if *receiverMode {
		l, err := receiver.NewReceiver(*address)
		if err != nil {
			log.Panic().Msg(err.Error())
		}

		l.Listen()
	}
}
