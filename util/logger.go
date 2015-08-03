package util

import log "github.com/Sirupsen/logrus"

var (
	Log *log.Logger
)

func init() {
	Log = log.New()
	Log.Level=log.DebugLevel
}
