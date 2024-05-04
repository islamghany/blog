INSERT INTO users (
        email,
        username,
        first_name,
        last_name,
        roles,
        password_hashed,
        enabled,
        created_at,
        updated_at
    )
VALUES (
        'user1@example.com',
        'user1',
        'John',
        'Doe',
        '{"user"}',
        'hashed_password_1',
        true,
        NOW(),
        NOW()
    ),
    (
        'user2@example.com',
        'user2',
        'Jane',
        'Smith',
        '{"user"}',
        'hashed_password_2',
        true,
        NOW(),
        NOW()
    ),
    (
        'admin@example.com',
        'admin',
        'Admin',
        'User',
        '{"admin"}',
        'hashed_admin_password',
        true,
        NOW(),
        NOW()
    );