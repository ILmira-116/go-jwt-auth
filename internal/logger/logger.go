package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	//Настройки логгера
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)

	//Формат вывода логов
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

}
