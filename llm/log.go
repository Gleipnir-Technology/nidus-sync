package llm

import (
	"log"
	"strings"

	"github.com/rs/zerolog"
	//"go.mau.fi/util/exzerolog"
)

type Logger = zerolog.Logger

func linkLogger(logger *zerolog.Logger) {
	//exzerolog.SetupDefaults(logger)
}

type ZerologWriter struct {
	zerologger zerolog.Logger
	level      zerolog.Level
}

func (w ZerologWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSuffix(string(p), "\n")
	event := w.zerologger.WithLevel(w.level)
	event.Msg(msg)
	return len(p), nil
}

func LoggerShim(l zerolog.Logger) *log.Logger {
	writer := &ZerologWriter{
		zerologger: l,
		level:      zerolog.DebugLevel,
	}
	return log.New(writer, "", 0)
}
