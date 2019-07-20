package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jaysoncena/bla/io/ratelimit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrNotPipe = errors.New("Not reading from a pipe")
)

const (
	defaultBufferSize = 1 << 16 // 64k
	// defaultMTU        = 1500
)

type Reader struct {
	bufferSize uint64
	buf        []byte
}

func NewReader() *Reader {
	return &Reader{
		bufferSize: defaultBufferSize,
		buf:        make([]byte, defaultBufferSize),
	}
}

func init() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Str("defaultBufferSize", strconv.Itoa(defaultBufferSize)).
		Msg("")
}

func (r *Reader) newReaderWithBuffer(src io.Reader) *bufio.Reader {
	return bufio.NewReaderSize(src, int(r.bufferSize))
}

func (r *Reader) ReadFromStdin() error {
	// Ticker
	ticker := time.NewTicker(time.Second)
	tick := make(chan time.Time, 1)

	go func() {
		for t := range ticker.C {
			tick <- t
		}
	}()
	// Put Now() so that it can be consumed immediately and
	// no need to wait for the tick
	tick <- time.Now()

	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Debug().
			Str("mode", v(info.Mode())).
			Str("size", v(info.Size())).
			Msg("")
		return ErrNotPipe
	}
	log.Debug().Msg("Pass pipe validation")

	// reader := bufio.NewReaderSize(os.Stdin, defaultBufferSize)
	reader := r.newReaderWithBuffer(os.Stdin)

	// for range tick {
	for {
		bytesRead, err := reader.Read(r.buf)
		log.Debug().Msgf("Read %d bytes", bytesRead)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// log.Debug().Msgf("Read %d bytes from stdin: %s", bytesRead, string(r.buf[:bytesRead]))
		os.Stdout.Write(r.buf[:bytesRead])
	}

	return nil
}

func (r *Reader) CopyWithRatelimit(dst io.Writer, src io.Reader, bps uint64) {
	reader := r.newReaderWithBuffer(src)
	limitedReader := ratelimit.NewReader(reader)
	limitedReader.SetRatelimit(bps)

	for {
		bytesRead, err := io.Copy(dst, limitedReader)
		if err == io.EOF || bytesRead == 0 {
			log.Info().
				Str("err", fmt.Sprintf("%+v", err)).
				Str("bytesRead", strconv.Itoa(int(bytesRead))).
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
