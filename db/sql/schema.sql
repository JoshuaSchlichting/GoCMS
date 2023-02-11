create table if not exists public.user
(
    id bigserial not null
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

create table if not exists public.organization
(
    id bigserial not null
        constraint organization_pk
            primary key,
    name text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists public.file
(
    id bigserial not null
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
    id bigserial not null
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
    id bigserial not null
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