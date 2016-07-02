package socketcrutch

import (
	logger "github.com/op/go-logging"
	"os"
	"strings"
)

var log *logger.Logger = logger.MustGetLogger("SocketCrutch")

// setup logger for package, noop by default
func init() {
	var format = logger.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backend := logger.NewLogBackend(os.Stderr, "SocketCrutch ", 0)
	backendFormatter := logger.NewBackendFormatter(backend, format)

	backendLeveled := logger.AddModuleLevel(backend)

	if strings.HasSuffix(os.Getenv("DEBUG"), "socketcrutch") {
		backendLeveled.SetLevel(logger.DEBUG, "")
	} else {
		backendLeveled.SetLevel(logger.CRITICAL, "")
	}

	logger.SetBackend(backendLeveled, backendFormatter)
}
