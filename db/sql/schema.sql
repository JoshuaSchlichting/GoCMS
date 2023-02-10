create table if not exists public.user
(
    id serial not null
        constraint user_pk
            primary key,
    name text not null unique,
    email text not null,
    attributes jsonb not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table if not exists public.file
(
    id serial not null
        constraint files_pk
            primary key,
    name text not null,
    blob bytea not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    owner_id integer not null
        constraint files_user_id_fk
            references public.user
);

create table if not exists public.messages
(
    id serial not null
        constraint messages_pk
            primary key,
    message text not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    owner_id integer not null
        constraint messages_user_id_fk
            references public.user
);