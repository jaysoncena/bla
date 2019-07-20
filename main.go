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
	flag         = kingpin.New("bla", "Multicast File Transfer")
	senderMode   = flag.Flag("sender", "Sender Mode").Bool()
	receiverMode = flag.Flag("receiver", "Receiver Mode").Bool()
	// address      = flag.Flag("addr", "IPv4 Address (Listen IP or Target IP)").Required().String()
	address  = flag.Flag("addr", "IPv4 Address (Listen IP or Target IP)").Required().String()
	testPipe = flag.Flag("testpipe", "Pipe Test").Bool()
	DebugLog = flag.Flag("debug", "Log level").Bool()
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *DebugLog {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	kingpin.MustParse(flag.Parse(os.Args[1:]))

	if !*senderMode && !*receiverMode && !*testPipe {
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
	} else if *testPipe {
		// err := r.ReadFromStdin()
		// if err != nil {
		// 	log.Error().Err(err)
		// }
		// r.ReadWithRatelimit(os.Stdout, os.Stdin, 1024*1024)

		// initialSpeed := 20

		// reader := ratelimit.NewReader(os.Stdin)
		// reader.SetRatelimit(20 * 1024 * 1024)

		// go func() {
		// 	ticker := time.NewTicker(60 * time.Second)

		// 	for range ticker.C {
		// 		log.Panic().Msg("Done!")
		// 		// if initialSpeed <= 10 {
		// 		// 	break
		// 		// }

		// 		// reader.UpdateLimit(uint64(initialSpeed * 1024 * 1024))

		// 		// initialSpeed -= 1
		// 	}

		// 	return
		// }()

		// for {
		// 	bytesRead, err := io.Copy(os.Stdout, reader)
		// 	if err == io.EOF || bytesRead == 0 {
		// 		log.Info().
		// 			Str("err", fmt.Sprintf("%+v", err)).
		// 			Str("bytesRead", strconv.Itoa(int(bytesRead))).
		// 			Msg("EOF for source")
		// 		break
		// 	} else if err != nil {
		// 		log.Error().Msgf("Error: %+v", err)
		// 		break
		// 	}
		// }

		r := io.NewReader()
		r.CopyWithRatelimit(os.Stdout, os.Stdin, 20<<20)

		log.Info().Msg("Done!")
		return

	}
}
