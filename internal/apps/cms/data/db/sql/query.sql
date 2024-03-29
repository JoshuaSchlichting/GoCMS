-- name: GetUser :one
SELECT * FROM public.user
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM public.user
WHERE name = $1;

-- name: ListUsers :many
SELECT * FROM public.user
ORDER BY name;

-- name: CreateUser :one
INSERT INTO public.user (
    id, name, email, attributes, created_ts, updated_ts
) VALUES (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
RETURNING *;

-- name: UpdateUser :one
UPDATE public.user
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM public.user
WHERE id = $1;

-- name: ListFiles :many
SELECT * FROM public.file
ORDER BY name;

-- name: UploadFile :one
INSERT INTO public.file (
    name, blob, created_ts, updated_ts, owner_id
) VALUES (
  $1, $2, current_timestamp, current_timestamp, $3
)
RETURNING *;

-- name: UpdateOrganization :one
update public.organization
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
WHERE id = $1
returning *;

-- name: CreateOrganization :one
insert into public.organization (
    id, name, email, attributes, created_ts, updated_ts
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
returning *;

-- name: DeleteOrganization :exec
delete from public.organization
where id = $1;

-- name: CreateUserGroup :one
insert into public.usergroup (
  id, name, email, attributes, created_ts, updated_ts
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp
)
returning *;

-- name: UpdateUserGroup :one
update public.usergroup
  set name = $2,
    email = $3,
    attributes = $4,
    updated_ts = current_timestamp
where id = $1
returning *;

-- name: DeleteUserGroup :exec
delete from public.usergroup
where id = $1;

-- name: ListUserGroups :many
select * from public.usergroup
order by name;

-- name: GetUserIsInGroup :one
select true
from
  public.user_usergroup
  left join public.usergroup
    on user_usergroup.usergroup_id = usergroup.id
where user_id = @user_id::uuid and usergroup.name = @usergroup_name::text;

-- name: ListOrganizations :many
SELECT * FROM public.organization
ORDER BY name;

-- name: ListMessagesTo :many
select
  id,
  from_id,
  subject,
  message,
  created_ts,
  updated_ts
from public.message
where to_username = $1;

-- name: ListMessagesFrom :many
select
  id,
  to_username,
  subject,
  message,
  created_ts,
  updated_ts
from public.message
where from_id = $1;

-- name: CreateMessage :one
insert into public.message (
  id, to_username, subject, message, created_ts, updated_ts, from_id
) values (
  $1, $2, $3, $4, current_timestamp, current_timestamp, $5
)
returning *;

-- name: CreateBlogPost :one
insert into public.blog_post (
  id, title, subtitle, featured_image_uri, body, author_id, created_ts, updated_ts
) values (
  $1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp
)
returning *;

-- name: ListBlogPosts :many
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
order by created_ts desc;

-- name: ListBlogPostsByUser :many
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
order by created_ts desc;

-- name: GetBlogPost :one
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
where id = $1;

-- name: SetMessageRead :exec
update public.message
  set read = $2
where id = $1;