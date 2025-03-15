package logger

import (
    "os"
    "github.com/rs/zerolog"
)

var log = zerolog.New(os.Stdout).With().Timestamp().Logger()

func GetLogger() zerolog.Logger {
    return log
}
