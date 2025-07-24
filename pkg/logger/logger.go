// Proje genelinde kullanılacak profesyonel logger.
// Logrus ile JSON formatında, seviyeli ve modüler loglama sağlar.

package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger global olarak kullanılacak logrus instance'ı
var Logger = logrus.New()

func init() {
	// Logları JSON formatında bas
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	// Log seviyesini ayarla (info, warn, error, debug, trace)
	Logger.SetLevel(logrus.InfoLevel)
	// Logları stdout'a yaz
	Logger.SetOutput(os.Stdout)
}

// Kullanım örneği:
// logger.Logger.Info("Kullanıcı giriş yaptı", logrus.Fields{"user_id": 123})
// logger.Logger.WithFields(logrus.Fields{"err": err}).Error("DB hatası")
