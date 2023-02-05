package data

import (
	"os"

	"github.com/joshuaschlichting/gocms/models"
)

type StubData struct{}

func (d StubData) GetUser(authUserId string) (models.User, error) {
	return models.User{
		UserName: "JohnDoe",
		Email:    "john.doe@email.com",
		Attributes: []string{
			"attribute1",
			"attribute2",
		},
	}, nil
}

func (d StubData) UploadFile(data []byte, fileName, userId string) error {
	// open file for writing and write data to it
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
