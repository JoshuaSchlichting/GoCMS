// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BlogPost struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	FeaturedImageURI string    `json:"featured_image_uri"`
	Body             string    `json:"body"`
	AuthorID         uuid.UUID `json:"author_id"`
	CreatedTS        time.Time `json:"created_ts"`
	UpdatedTS        time.Time `json:"updated_ts"`
}

type File struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Blob      []byte    `json:"blob"`
	CreatedTS time.Time `json:"created_ts"`
	UpdatedTS time.Time `json:"updated_ts"`
	OwnerID   uuid.UUID `json:"owner_id"`
}

type FileFilegroup struct {
	FileID      uuid.UUID `json:"file_id"`
	FilegroupID uuid.UUID `json:"filegroup_id"`
	CreatedTS   time.Time `json:"created_ts"`
	UpdatedTS   time.Time `json:"updated_ts"`
}

type Filegroup struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedTS time.Time `json:"created_ts"`
	UpdatedTS time.Time `json:"updated_ts"`
}

type Invoice struct {
	ID             uuid.UUID `json:"id"`
	Amount         float64   `json:"amount"`
	CreatedTS      time.Time `json:"created_ts"`
	UpdatedTS      time.Time `json:"updated_ts"`
	UserID         uuid.UUID `json:"user_id"`
	OrgnaizationID uuid.UUID `json:"orgnaization_id"`
}

type Message struct {
	ID         uuid.UUID `json:"id"`
	ToUsername string    `json:"to_username"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	CreatedTS  time.Time `json:"created_ts"`
	UpdatedTS  time.Time `json:"updated_ts"`
	FromID     uuid.UUID `json:"from_id"`
}

type Organization struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
	CreatedTS  time.Time       `json:"created_ts"`
	UpdatedTS  time.Time       `json:"updated_ts"`
}

type PermissionAttribute struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedTS time.Time `json:"created_ts"`
	UpdatedTS time.Time `json:"updated_ts"`
}

type User struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.NullUUID   `json:"organization_id"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Attributes     json.RawMessage `json:"attributes"`
	CreatedTS      time.Time       `json:"created_ts"`
	UpdatedTS      time.Time       `json:"updated_ts"`
}

type UserPermissionAttribute struct {
	UserID                uuid.UUID `json:"user_id"`
	PermissionAttributeID uuid.UUID `json:"permission_attribute_id"`
	CreatedTS             time.Time `json:"created_ts"`
	UpdatedTS             time.Time `json:"updated_ts"`
}

type UserUsergroup struct {
	UserID      uuid.UUID `json:"user_id"`
	UsergroupID uuid.UUID `json:"usergroup_id"`
	CreatedTS   time.Time `json:"created_ts"`
	UpdatedTS   time.Time `json:"updated_ts"`
}

type Usergroup struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
	CreatedTS  time.Time       `json:"created_ts"`
	UpdatedTS  time.Time       `json:"updated_ts"`
}

type UsergroupOrganization struct {
	UsergroupID    uuid.UUID `json:"usergroup_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedTS      time.Time `json:"created_ts"`
	UpdatedTS      time.Time `json:"updated_ts"`
}

type UsergroupPermissionAttribute struct {
	UsergroupID           uuid.UUID `json:"usergroup_id"`
	PermissionAttributeID uuid.UUID `json:"permission_attribute_id"`
	CreatedTS             time.Time `json:"created_ts"`
	UpdatedTS             time.Time `json:"updated_ts"`
}