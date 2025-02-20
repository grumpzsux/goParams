package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Init configures the global logger.
func Init(verbose bool) {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}
