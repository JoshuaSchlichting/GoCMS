// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const createBlogPost = `-- name: CreateBlogPost :one
insert into public.blog_post (
  id, title, subtitle, featured_image_uri, body, author_id, created_ts, updated_ts
) values (
  $1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp
)
returning id, title, subtitle, featured_image_uri, body, author_id, created_ts, updated_ts
`

type CreateBlogPostParams struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	FeaturedImageURI string    `json:"featured_image_uri"`
	Body             string    `json:"body"`
	AuthorID         uuid.UUID `json:"author_id"`
}

func (q *Queries) CreateBlogPost(ctx context.Context, arg CreateBlogPostParams) (BlogPost, error) {
	row := q.db.QueryRowContext(ctx, createBlogPost,
		arg.ID,
		arg.Title,
		arg.Subtitle,
		arg.FeaturedImageURI,
		arg.Body,
		arg.AuthorID,
	)
	var i BlogPost
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Subtitle,
		&i.FeaturedImageURI,
		&i.Body,
		&i.AuthorID,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const createMessage = `-- name: CreateMessage :one
insert into public.message (
  id, to_username, subject, message, created_ts, updated_ts, from_id
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp, $5
)
returning id, to_username, subject, message, read, created_ts, updated_ts, from_id
`

type CreateMessageParams struct {
	ID         uuid.UUID `json:"id"`
	ToUsername string    `json:"to_username"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	FromID     uuid.UUID `json:"from_id"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.ToUsername,
		arg.Subject,
		arg.Message,
		arg.FromID,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.ToUsername,
		&i.Subject,
		&i.Message,
		&i.Read,
		&i.CreatedTS,
		&i.UpdatedTS,
		&i.FromID,
	)
	return i, err
}

const createOrganization = `-- name: CreateOrganization :one
insert into public.organization (
    id, name, email, attributes, created_ts, updated_ts
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
returning id, name, email, attributes, created_ts, updated_ts
`

type CreateOrganizationParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (Organization, error) {
	row := q.db.QueryRowContext(ctx, createOrganization,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO public.user (
    id, name, email, attributes, created_ts, updated_ts
) VALUES (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
RETURNING id, organization_id, name, email, attributes, created_ts, updated_ts
`

type CreateUserParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const createUserGroup = `-- name: CreateUserGroup :one
insert into public.usergroup (
  id, name, email, attributes, created_ts, updated_ts
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
returning id, name, email, attributes, created_ts, updated_ts
`

type CreateUserGroupParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) CreateUserGroup(ctx context.Context, arg CreateUserGroupParams) (Usergroup, error) {
	row := q.db.QueryRowContext(ctx, createUserGroup,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i Usergroup
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const deleteOrganization = `-- name: DeleteOrganization :exec
delete from public.organization
where id = $1
`

func (q *Queries) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteOrganization, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM public.user
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const deleteUserGroup = `-- name: DeleteUserGroup :exec
delete from public.usergroup
where id = $1
`

func (q *Queries) DeleteUserGroup(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUserGroup, id)
	return err
}

const getBlogPost = `-- name: GetBlogPost :one
select
  id,
  title,
  subtitle,
  body,
  author_id,
  featured_image_uri,
  created_ts,
  updated_ts
from public.blog_post
where id = $1
`

type GetBlogPostRow struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Body             string    `json:"body"`
	AuthorID         uuid.UUID `json:"author_id"`
	FeaturedImageURI string    `json:"featured_image_uri"`
	CreatedTS        time.Time `json:"created_ts"`
	UpdatedTS        time.Time `json:"updated_ts"`
}

func (q *Queries) GetBlogPost(ctx context.Context, id uuid.UUID) (GetBlogPostRow, error) {
	row := q.db.QueryRowContext(ctx, getBlogPost, id)
	var i GetBlogPostRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Subtitle,
		&i.Body,
		&i.AuthorID,
		&i.FeaturedImageURI,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, organization_id, name, email, attributes, created_ts, updated_ts FROM public.user
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, organization_id, name, email, attributes, created_ts, updated_ts FROM public.user
WHERE name = $1
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const getUserIsInGroup = `-- name: GetUserIsInGroup :one
select true
from
  public.user_usergroup
  left join public.usergroup
    on user_usergroup.usergroup_id = usergroup.id
where user_id = $1::uuid and usergroup.name = $2::text
`

type GetUserIsInGroupParams struct {
	UserID        uuid.UUID `json:"user_id"`
	UsergroupName string    `json:"usergroup_name"`
}

func (q *Queries) GetUserIsInGroup(ctx context.Context, arg GetUserIsInGroupParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, getUserIsInGroup, arg.UserID, arg.UsergroupName)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const listBlogPosts = `-- name: ListBlogPosts :many
select
  id,
  title,
  subtitle,
  body,
  author_id,
  featured_image_uri,
  created_ts,
  updated_ts
from public.blog_post
order by created_ts desc
`

type ListBlogPostsRow struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Body             string    `json:"body"`
	AuthorID         uuid.UUID `json:"author_id"`
	FeaturedImageURI string    `json:"featured_image_uri"`
	CreatedTS        time.Time `json:"created_ts"`
	UpdatedTS        time.Time `json:"updated_ts"`
}

func (q *Queries) ListBlogPosts(ctx context.Context) ([]ListBlogPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, listBlogPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBlogPostsRow
	for rows.Next() {
		var i ListBlogPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Subtitle,
			&i.Body,
			&i.AuthorID,
			&i.FeaturedImageURI,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBlogPostsByUser = `-- name: ListBlogPostsByUser :many
select
  id,
  title,
  subtitle,
  body,
  author_id,
  featured_image_uri,
  created_ts,
  updated_ts
from public.blog_post
where author_id = $1
order by created_ts desc
`

type ListBlogPostsByUserRow struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Body             string    `json:"body"`
	AuthorID         uuid.UUID `json:"author_id"`
	FeaturedImageURI string    `json:"featured_image_uri"`
	CreatedTS        time.Time `json:"created_ts"`
	UpdatedTS        time.Time `json:"updated_ts"`
}

func (q *Queries) ListBlogPostsByUser(ctx context.Context, authorID uuid.UUID) ([]ListBlogPostsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, listBlogPostsByUser, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBlogPostsByUserRow
	for rows.Next() {
		var i ListBlogPostsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Subtitle,
			&i.Body,
			&i.AuthorID,
			&i.FeaturedImageURI,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFiles = `-- name: ListFiles :many
SELECT id, name, blob, created_ts, updated_ts, owner_id FROM public.file
ORDER BY name
`

func (q *Queries) ListFiles(ctx context.Context) ([]File, error) {
	rows, err := q.db.QueryContext(ctx, listFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Blob,
			&i.CreatedTS,
			&i.UpdatedTS,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMessagesFrom = `-- name: ListMessagesFrom :many
select
  id,
  to_username,
  subject,
  message,
  created_ts,
  updated_ts
from public.message
where from_id = $1
`

type ListMessagesFromRow struct {
	ID         uuid.UUID `json:"id"`
	ToUsername string    `json:"to_username"`
	Subject    string    `json:"subject"`
	Message    string    `json:"message"`
	CreatedTS  time.Time `json:"created_ts"`
	UpdatedTS  time.Time `json:"updated_ts"`
}

func (q *Queries) ListMessagesFrom(ctx context.Context, fromID uuid.UUID) ([]ListMessagesFromRow, error) {
	rows, err := q.db.QueryContext(ctx, listMessagesFrom, fromID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListMessagesFromRow
	for rows.Next() {
		var i ListMessagesFromRow
		if err := rows.Scan(
			&i.ID,
			&i.ToUsername,
			&i.Subject,
			&i.Message,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMessagesTo = `-- name: ListMessagesTo :many
select
  id,
  from_id,
  subject,
  message,
  created_ts,
  updated_ts
from public.message
where to_username = $1
`

type ListMessagesToRow struct {
	ID        uuid.UUID `json:"id"`
	FromID    uuid.UUID `json:"from_id"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedTS time.Time `json:"created_ts"`
	UpdatedTS time.Time `json:"updated_ts"`
}

func (q *Queries) ListMessagesTo(ctx context.Context, toUsername string) ([]ListMessagesToRow, error) {
	rows, err := q.db.QueryContext(ctx, listMessagesTo, toUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListMessagesToRow
	for rows.Next() {
		var i ListMessagesToRow
		if err := rows.Scan(
			&i.ID,
			&i.FromID,
			&i.Subject,
			&i.Message,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOrganizations = `-- name: ListOrganizations :many
SELECT id, name, email, attributes, created_ts, updated_ts FROM public.organization
ORDER BY name
`

func (q *Queries) ListOrganizations(ctx context.Context) ([]Organization, error) {
	rows, err := q.db.QueryContext(ctx, listOrganizations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Organization
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Attributes,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserGroups = `-- name: ListUserGroups :many
select id, name, email, attributes, created_ts, updated_ts from public.usergroup
order by name
`

func (q *Queries) ListUserGroups(ctx context.Context) ([]Usergroup, error) {
	rows, err := q.db.QueryContext(ctx, listUserGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Usergroup
	for rows.Next() {
		var i Usergroup
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Attributes,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, organization_id, name, email, attributes, created_ts, updated_ts FROM public.user
ORDER BY name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.OrganizationID,
			&i.Name,
			&i.Email,
			&i.Attributes,
			&i.CreatedTS,
			&i.UpdatedTS,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setMessageRead = `-- name: SetMessageRead :exec
update public.message
  set read = $2
where id = $1
`

type SetMessageReadParams struct {
	ID   uuid.UUID `json:"id"`
	Read bool      `json:"read"`
}

func (q *Queries) SetMessageRead(ctx context.Context, arg SetMessageReadParams) error {
	_, err := q.db.ExecContext(ctx, setMessageRead, arg.ID, arg.Read)
	return err
}

const updateOrganization = `-- name: UpdateOrganization :one
update public.organization
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
WHERE id = $1
returning id, name, email, attributes, created_ts, updated_ts
`

type UpdateOrganizationParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (Organization, error) {
	row := q.db.QueryRowContext(ctx, updateOrganization,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i Organization
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE public.user
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
WHERE id = $1
RETURNING id, organization_id, name, email, attributes, created_ts, updated_ts
`

type UpdateUserParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OrganizationID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const updateUserGroup = `-- name: UpdateUserGroup :one
update public.usergroup
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
where id = $1
returning id, name, email, attributes, created_ts, updated_ts
`

type UpdateUserGroupParams struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Attributes json.RawMessage `json:"attributes"`
}

func (q *Queries) UpdateUserGroup(ctx context.Context, arg UpdateUserGroupParams) (Usergroup, error) {
	row := q.db.QueryRowContext(ctx, updateUserGroup,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Attributes,
	)
	var i Usergroup
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Attributes,
		&i.CreatedTS,
		&i.UpdatedTS,
	)
	return i, err
}

const uploadFile = `-- name: UploadFile :one
INSERT INTO public.file (
    name, blob, created_ts, updated_ts, owner_id
) VALUES (
  $1, $2, current_timestamp, current_timestamp, $3
)
RETURNING id, name, blob, created_ts, updated_ts, owner_id
`

type UploadFileParams struct {
	Name    string    `json:"name"`
	Blob    []byte    `json:"blob"`
	OwnerID uuid.UUID `json:"owner_id"`
}

func (q *Queries) UploadFile(ctx context.Context, arg UploadFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, uploadFile, arg.Name, arg.Blob, arg.OwnerID)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Blob,
		&i.CreatedTS,
		&i.UpdatedTS,
		&i.OwnerID,
	)
	return i, err
}
