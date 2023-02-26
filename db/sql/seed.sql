-- seed users for
-- create table if not exists public.user
-- (
--     id bigserial not null
--         constraint user_pk
--             primary key,
--     organization_id integer
--         constraint user_organization_id_fk
--             references public.organization,
--     name text not null unique,
--     email text not null,
--     attributes jsonb not null,
--     created_at timestamp not null default current_timestamp,
--     updated_at timestamp not null default current_timestamp
-- );

insert into public.user (name, email, attributes, created_at, updated_at)
values ('admin',
        'admin@localhost',
        '{"role": "admin"}',
        current_timestamp,
        current_timestamp),
        (
            'user',
            'user@localhost',
            '{"role": "user"}',
            current_timestamp,
            current_timestamp
        ),
        (
            'guest',
            'guest@localhost',
            '{"role": "guest"}',
            current_timestamp,
            current_timestamp
        ),
        (
            'test',
            'test@localhost',
            '{"role": "test"}',
            current_timestamp,
            current_timestamp
        ),
        (
            'test2',
            'test2@localhost',
            '{"role": "test2"}',
            current_timestamp,
            current_timestamp
        ),
        (
            'test3',
            'test3@localhost',
            '{"role": "test3"}',
            current_timestamp,
            current_timestamp
        )
    ;

-- seed public.permission_attribute with permissions such as 'create user', 'delete user', 'create organization', 'delete organization', etc.
insert into public.permission_attribute (name, created_at, updated_at)
values
    ('is superuser',
        current_timestamp,
        current_timestamp
    ),
    ('create user',
        current_timestamp,
        current_timestamp
    ),
    (
        'delete user',
        current_timestamp,
        current_timestamp
    ),
    (
        'create organization',
        current_timestamp,
        current_timestamp
    ),
    (
        'delete organization',
        current_timestamp,
        current_timestamp
    ),
    (
        'create group',
        current_timestamp,
        current_timestamp
    ),
    (
        'delete group',
        current_timestamp,
        current_timestamp
    ),
    (
        'create file',
        current_timestamp,
        current_timestamp
    ),
    (
        'delete file',
        current_timestamp,
        current_timestamp
    ),
    (
        'create permission',
        current_timestamp,
        current_timestamp
    ),
    (
        'delete permission',
        current_timestamp,
        current_timestamp
    )
;