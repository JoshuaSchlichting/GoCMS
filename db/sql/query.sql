-- name: GetUser :one
SELECT * FROM public.user
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM public.user
WHERE name = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM public.user
ORDER BY name;

-- name: CreateUser :one
INSERT INTO public.user (
    name, email, attributes, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUser :one
UPDATE public.user
  set name = $2,
    email = $3,
    attributes = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM public.user
WHERE id = $1;

-- name: ListFiles :many
SELECT * FROM public.file
ORDER BY name;

-- name: UpdatUser :one
UPDATE public.user
  set name = $2,
    email = $3,
    attributes = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: UploadFile :one
INSERT INTO public.file (
    name, blob, created_at, updated_at, owner_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

