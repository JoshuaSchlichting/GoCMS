package data

import (
	"github.com/joshuaschlichting/gocms/models"
)

type Data interface {
	GetUser(authUserId string) (models.User, error)
}

type StubData struct {
}

func (d *StubData) GetUser(authUserId string) (models.User, error) {
	return models.User{
		UserName: "JohnDoe",
		Email:    "john.doe@email.com",
		Attributes: []string{
			"attribute1",
			"attribute2",
		},
	}, nil
}
