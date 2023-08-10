insert into public.user (id, name, email, attributes, created_ts, updated_ts)
values 
    ('5a0e233a-6b11-4d47-8820-1f0b8ca630a2', 'admin', 'admin@localhost', '{"role": "admin"}', current_timestamp, current_timestamp),
    ('1e0f875a-2ab3-4d9f-8b5a-4e4d4b4f4b7c', 'user', 'user@localhost', '{"role": "user"}', current_timestamp, current_timestamp),
    ('d846fcbf-2dc2-4d31-99d9-18b9f1faa138', 'guest', 'guest@localhost', '{"role": "guest"}', current_timestamp, current_timestamp),
    ('1dfe9b46-ecb8-478d-9d11-68c2b06dcf77', 'test', 'test@localhost', '{"role": "test"}', current_timestamp, current_timestamp),
    ('3ed55756-25e8-44e9-9d67-c17a68de94e5', 'test2', 'test2@localhost', '{"role": "test2"}', current_timestamp, current_timestamp),
    ('ca7e719e-8d4f-4f26-84b7-0f6a28f7d193', 'test3', 'test3@localhost', '{"role": "test3"}', current_timestamp, current_timestamp);


-- seed public.permission_attribute with permissions such as 'create user', 'delete user', 'create organization', 'delete organization', etc.
insert into public.permission_attribute (name, created_ts, updated_ts)
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

INSERT INTO public.message (id, to_id, subject, message, created_ts, updated_ts, from_id)
VALUES
('123e4567-e89b-12d3-a456-426655440000', '123e4567-e89b-12d3-a456-426655440001', 'Sample Subject 1', 'This is a sample message 1', current_timestamp, current_timestamp, '123e4567-e89b-12d3-a456-426655440002'),
('123e4567-e89b-12d3-a456-426655440003', '123e4567-e89b-12d3-a456-426655440004', 'Sample Subject 2', 'This is a sample message 2', current_timestamp, current_timestamp, '123e4567-e89b-12d3-a456-426655440005'),
('123e4567-e89b-12d3-a456-426655440006', '123e4567-e89b-12d3-a456-426655440007', 'Sample Subject 3', 'This is a sample message 3', current_timestamp, current_timestamp, '123e4567-e89b-12d3-a456-426655440008');