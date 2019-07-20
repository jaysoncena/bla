package ratelimit

import (
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

const (
	defaultBurst = 1 << 16 // 64KiB
	defaultBps   = 1 << 20 // 1MiB
)

type Reader struct {
	bps     uint64
	burst   uint64
	limiter *rate.Limiter
	r       io.Reader
	ctx     context.Context
}

type Writer struct {
	bps     uint64
	burst   uint64
	limiter *rate.Limiter
	w       io.Writer
	ctx     context.Context
}

func init() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().
		Str("bps", strconv.Itoa(defaultBps)).
		Str("burst", strconv.Itoa(defaultBurst)).
		Msg("")
}

func NewReader(r io.Reader) *Reader {
	log.Debug().
		Str("bps", strconv.Itoa(defaultBps)).
		Str("burst", strconv.Itoa(defaultBurst)).
		Msg("Reader with ratelimit created")
	return &Reader{
		r:     r,
		bps:   defaultBps,
		burst: defaultBurst,
		ctx:   context.Background(),
	}
}

func NewWriter(w io.Writer) *Writer {
	log.Debug().
		Str("bps", string(defaultBps)).
		Str("burst", string(defaultBurst)).
		Msg("Reader with ratelimit created")
	return &Writer{
		w:     w,
		bps:   defaultBps,
		burst: defaultBurst,
		ctx:   context.Background(),
	}
}

func (r *Reader) SetRatelimit(bps uint64) {
	r.bps = bps

	r.limiter = rate.NewLimiter(rate.Limit(bps), int(r.burst))
	r.limiter.AllowN(time.Now(), int(r.burst))
	log.Debug().
		Str("bps", strconv.Itoa(int(r.bps))).
		Str("burst", strconv.Itoa(int(r.burst))).
		Msg("SetRatelimit()")
}

func (r *Reader) UpdateLimit(bps uint64) {
	r.bps = bps
	r.limiter.SetLimit(rate.Limit(bps))
	log.Debug().
		Str("bps", string(bps)).
		Msg("Read limit updated")
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.limiter == nil {
		return r.r.Read(p)
	}

	n, err := r.r.Read(p)
	if err != nil {
		return n, err
	}

	if err := r.limiter.WaitN(r.ctx, n); err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) SetRatelimit(bps uint64) {
	w.bps = bps

	w.limiter = rate.NewLimiter(rate.Limit(bps), int(w.burst))
	w.limiter.AllowN(time.Now(), int(w.burst))
	log.Debug().
		Str("bps", strconv.Itoa(int(w.bps))).
		Str("burst", strconv.Itoa(int(w.burst))).
		Msg("SetRatelimit()")
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.limiter == nil {
		return w.w.Write(p)
	}

	n, err := w.w.Write(p)
	if err != nil {
		return n, err
	}

	if err := w.limiter.WaitN(w.ctx, n); err != nil {
		return n, err
	}

	return n, err
}
