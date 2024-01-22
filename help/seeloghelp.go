package help

import "github.com/cihub/seelog"

func SetupLogger() {
	logger, err := seelog.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		return
	}
	seelog.ReplaceLogger(logger)
}
