// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptorLog

import (
	"github.com/rs/zerolog/log"
)

func Trace(format string, v ...interface{}) {
	log.Trace().Msgf(format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Err(err error, format string, v ...interface{}) {
	log.Err(err).Msgf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Panic(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}
