package data

import (
	"database/sql"

	"github.com/joshuaschlichting/gocms/models"
)

var db *sql.DB

func Init(dbConn *sql.DB) {
	db = dbConn
}

type Data interface {
	GetUser(authUserId string) (models.User, error)
	UploadFile(data []byte, fileName, userId string) error
}
