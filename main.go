package main

import (
	"bufio"
	"fmt"
	// "io"
	"os"
	// "time"
	// "strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jaysoncena/bla/io"
	// "github.com/jaysoncena/bla/io/ratelimit"
	"github.com/jaysoncena/bla/receiver"
	"github.com/jaysoncena/bla/sender"
)

var (
	senderCmd    = kingpin.Command("sender", "Sender Mode")
	senderAddr   = senderCmd.Flag("addr", "Target Multicast IP Address").Required().String()
	receiverCmd  = kingpin.Command("receiver", "Receiver Mode")
	receiverAddr = receiverCmd.Flag("addr", "Multicast IP address to listen to").Required().String()
	pipeCmd      = kingpin.Command("pipe", "Test command")
	debugLog     = kingpin.Flag("debug", "Debug Mode").Bool()

	// flag         = kingpin.New("bla", "Multicast File Transfer")
	// senderMode   = flag.Flag("sender", "Sender Mode").Bool()
	// receiverMode = flag.Flag("receiver", "Receiver Mode").Bool()
	// address      = flag.Flag("addr", "IPv4 Address (Listen IP or Target IP)").Required().String()
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugLog {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// kingpin.MustParse(flag.Parse(os.Args[1:]))
}

func main() {
	switch kingpin.Parse() {
	case senderCmd.FullCommand():
		log.Info().Msg("Running in SENDER mode")
		s, err := sender.NewSender(*senderAddr)
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
	case receiverCmd.FullCommand():
		log.Info().Msg("Running in RECEIVER mode")
		l, err := receiver.NewReceiver(*receiverAddr)
		if err != nil {
			log.Panic().Msg(err.Error())
		}

		l.Listen()
	case pipeCmd.FullCommand():
		io.CopyWithRatelimit(os.Stdout, os.Stdin, 1<<28)

		log.Info().Msg("Done!")
	}
}
