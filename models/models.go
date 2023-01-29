package models

type User struct {
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Attributes []string `json:"attributes"`
}
