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
        'user22@example.com',
        'islamghany',
        'Islam',
        'Ghany',
        '{"user"}',
        '$2a$10$nI9Ery433WBZ1VwqeAvvd.M/Jy4MVbVClgT1M/HzlV4TEy5DOWXZS',
        true,
        NOW(),
        NOW()
    );