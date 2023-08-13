create database cms;

create table if not exists public.organization
(
    id uuid not null primary key,
    name text not null,
    email text not null,
    attributes jsonb not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);

create table if not exists public.user
(
    id uuid not null primary key,
    organization_id uuid,
    name text not null unique,
    email text not null,
    attributes jsonb not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);

create table if not exists public.file
(
    id uuid not null
        constraint files_pk
            primary key,
    name text not null,
    blob bytea not null,
    created_ts timestamp not null  default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    owner_id uuid not null
        constraint files_user_id_fk
            references public.user
);

create table if not exists public.message
(
    id uuid not null primary key,
    to_username text not null,
    subject text not null,
    message text not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    from_id uuid not null
);

create table if not exists public.invoice
(
    id uuid not null primary key,
    amount float not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    user_id uuid not null,
    orgnaization_id uuid not null
);

create table if not exists public.usergroup
(
    id uuid not null primary key,
    name text not null,
    email text not null,
    attributes jsonb not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);

create table if not exists public.user_usergroup
(
    user_id uuid not null,
    usergroup_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    primary key (user_id, usergroup_id)
);

create table if not exists public.usergroup_organization
(
    usergroup_id uuid not null,
    organization_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    primary key (usergroup_id, organization_id)
);

create table if not exists public.permission_attribute
(
    id uuid not null primary key,
    name text not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);

create table if not exists public.usergroup_permission_attribute
(
    usergroup_id uuid not null,
    permission_attribute_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    primary key (usergroup_id, permission_attribute_id)
);

create table if not exists public.user_permission_attribute
(
    user_id uuid not null,
    permission_attribute_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    primary key (user_id, permission_attribute_id)
);

create table if not exists public.filegroup
(
    id uuid not null primary key,
    name text not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);

create table if not exists public.file_filegroup
(
    file_id uuid not null,
    filegroup_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp,
    primary key (file_id, filegroup_id)
);

create table blog_post
(
    id uuid not null primary key,
    title text not null,
    subtitle text not null,
    featured_image_uri text not null,
    body text not null,
    author_id uuid not null,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp
);