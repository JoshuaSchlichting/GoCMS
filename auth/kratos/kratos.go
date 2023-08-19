package kratos

import (
	"log/slog"

	"github.com/joshuaschlichting/gocms/config"
)

var conf *config.Config

var logger *slog.Logger

var kratos_whoami_url string

func InitKratos(c *config.Config) {
	conf = c
	kratos_whoami_url = "http://" + conf.Auth.Kratos.Host + ":" + conf.Auth.Kratos.Port + "/sessions/whoami"
}
func SetLogger(l *slog.Logger) {
	logger = l
}
