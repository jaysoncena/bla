package io

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jaysoncena/bla/io/ratelimit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultCopyBufferSize = 1 << 16
)

func init() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Str("defaultCopyBufferSize", strconv.Itoa(defaultCopyBufferSize)).
		Msg("")
}

func CopyWithRatelimit(dst io.Writer, src io.Reader, bps uint64) {
	writer := ratelimit.NewWriter(dst)
	if bps >= 0 {
		writer.SetRatelimit(bps)
	}

	buf := make([]byte, defaultCopyBufferSize)
	for {
		bytesWritten, err := io.CopyBuffer(writer, src, buf)
		if err != io.EOF || bytesWritten == 0 {
			log.Info().
				Str("err", v(err)).
				Str("bytesWritten", strconv.Itoa(int(bytesWritten))).
				Msg("EOF for source")
			break
		} else if err != nil {
			log.Error().Msgf("Error: %+v", err)
			break
		}
	}

}

func v(any interface{}) string {
	return fmt.Sprintf("%+v", any)
}
