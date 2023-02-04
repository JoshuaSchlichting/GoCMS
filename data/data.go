package data

import (
	"github.com/joshuaschlichting/gocms/models"
)

type Data interface {
	GetUser(authUserId string) (models.User, error)
}
