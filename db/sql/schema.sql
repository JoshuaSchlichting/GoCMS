create table public.user
(
    id serial not null
        constraint user_pk
            primary key,
    name text not null,
    email text not null,
    attributes jsonb not null,
    password text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table public.files
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