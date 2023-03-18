
create table if not exists public.organization
(
    id uuid not null
        constraint organization_pk
            primary key,
    name text not null,
    email text not null,
    attributes jsonb not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.user
(
    id uuid not null
        constraint user_pk
            primary key,
    organization_id integer
        constraint user_organization_id_fk
            references public.organization,
    name text not null unique,
    email text not null,
    attributes jsonb not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.file
(
    id uuid not null
        constraint files_pk
            primary key,
    name text not null,
    blob bytea not null,
    created_at timestamp not null  default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    owner_id integer not null
        constraint files_user_id_fk
            references public.user
);

create table if not exists public.message
(
    id uuid not null
        constraint messages_pk
            primary key,
    to_id integer not null
        constraint messages_to_fk
            references public.user,
    message text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    from_id integer not null
        constraint messages_user_id_fk
            references public.user
);

create table if not exists public.invoice
(
    id uuid not null
        constraint invoice_pk
            primary key,
    amount float not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    user_id integer not null
        constraint invoice_user_id_fk
            references public.user,
    orgnaization_id integer not null
        constraint invoice_orgnaization_id_fk
            references public.organization
);

create table if not exists public.usergroup
(
    id uuid not null
        constraint usergroup_pk
            primary key,
    name text not null,
    email text not null,
    attributes jsonb not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.user_usergroup
(
    user_id integer not null
        constraint user_usergroup_user_id_fk
            references public.user,
    usergroup_id integer not null
        constraint user_usergroup_usergroup_id_fk
            references public.usergroup,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    constraint user_usergroup_pk
        primary key (user_id, usergroup_id)
);

create table if not exists public.usergroup_organization
(
    usergroup_id integer not null
        constraint usergroup_organization_usergroup_id_fk
            references public.usergroup,
    organization_id integer not null
        constraint usergroup_organization_organization_id_fk
            references public.organization,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    constraint usergroup_organization_pk
        primary key (usergroup_id, organization_id)
);

create table if not exists public.permission_attribute
(
    id uuid not null
        constraint permission_attribute_pk
            primary key,
    name text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.usergroup_permission_attribute
(
    usergroup_id integer not null
        constraint usergroup_permission_attribute_usergroup_id_fk
            references public.usergroup,
    permission_attribute_id integer not null
        constraint usergroup_permission_attribute_permission_attribute_id_fk
            references public.permission_attribute,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    constraint usergroup_permission_attribute_pk
        primary key (usergroup_id, permission_attribute_id)
);

create table if not exists public.user_permission_attribute
(
    user_id integer not null
        constraint user_permission_attribute_user_id_fk
            references public.user,
    permission_attribute_id integer not null
        constraint user_permission_attribute_permission_attribute_id_fk
            references public.permission_attribute,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    constraint user_permission_attribute_pk
        primary key (user_id, permission_attribute_id)
);

create table if not exists public.filegroup
(
    id uuid not null
        constraint file_group_pk
            primary key,
    name text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.file_filegroup
(
    file_id integer not null
        constraint file_filegroup_file_id_fk
            references public.file,
    filegroup_id integer not null
        constraint file_filegroup_filegroup_id_fk
            references public.filegroup,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    constraint file_filegroup_pk
        primary key (file_id, filegroup_id)
);